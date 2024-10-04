package main

import (
	"github.com/cantylv/authorization-service/microservices/task_manager/config"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/app"
	"go.uber.org/zap"
)

func main() {
	logger := zap.Must(zap.NewProduction())
	config.Read("./config/config.yaml", logger)
	app.Run(logger)
}
