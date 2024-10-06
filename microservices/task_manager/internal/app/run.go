package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/cantylv/authorization-service/microservices/task_manager/internal/clients"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/delivery/route"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Run(logger *zap.Logger) {
	// создадим кластер клиентов наших микросервисов
	clientCluster := clients.InitCluster()

	r := mux.NewRouter()
	// инициализуруем серверные ручки
	handler := route.InitHTTPHandlers(r, clientCluster, logger)
	srv := &http.Server{
		Handler:      handler,
		Addr:         viper.GetString("task_manager.address"),
		WriteTimeout: viper.GetDuration("task_manager.write_timeout"),
		ReadTimeout:  viper.GetDuration("task_manager.read_timeout"),
		IdleTimeout:  viper.GetDuration("task_manager.idle_timeout"),
	}

	go func() {
		logger.Info(fmt.Sprintf("server has started at the address %s", viper.GetString("task_manager.address")))
		if err := srv.ListenAndServe(); err != nil {
			logger.Warn(fmt.Sprintf("error after end of receiving requests: %v", err))
		}
	}()

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("task_manager.shutdown_duration"))
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("server has shut down with an error: %v", err))
		os.Exit(1)
	}
	logger.Info("server has shut down")
}
