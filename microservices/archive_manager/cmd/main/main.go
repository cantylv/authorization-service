package main

import (
	"github.com/cantylv/authorization-service/microservices/archive_manager/config"
	"github.com/cantylv/authorization-service/microservices/archive_manager/internal/app"
	"go.uber.org/zap"
)

func main() {
	logger := zap.Must(zap.NewProduction())
	config.Read("./config/config.yaml", logger)
	app.Run(logger)
}
