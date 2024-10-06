package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/cantylv/authorization-service/microservices/archive_manager/internal/delivery/route"
	"github.com/cantylv/authorization-service/microservices/archive_manager/services/postgres"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Run(logger *zap.Logger) {
	// init psql
	postgresClient := postgres.Init(logger)
	r := mux.NewRouter()
	// инициализуруем серверные ручки
	handler := route.InitHTTPHandlers(r, postgresClient, logger)
	srv := &http.Server{
		Handler:      handler,
		Addr:         viper.GetString("archive_manager.address"),
		WriteTimeout: viper.GetDuration("archive_manager.write_timeout"),
		ReadTimeout:  viper.GetDuration("archive_manager.read_timeout"),
		IdleTimeout:  viper.GetDuration("archive_manager.idle_timeout"),
	}

	go func() {
		logger.Info(fmt.Sprintf("server has started at the address %s", viper.GetString("archive_manager.address")))
		if err := srv.ListenAndServe(); err != nil {
			logger.Warn(fmt.Sprintf("error after end of receiving requests: %v", err))
		}
	}()

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("archive_manager.shutdown_duration"))
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("server has shut down with an error: %v", err))
		os.Exit(1)
	}
	logger.Info("server has shut down")
}
