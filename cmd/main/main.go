package main

import (
	"github.com/cantylv/authorization-service/config"
	"github.com/cantylv/authorization-service/internal/app"
	"go.uber.org/zap"
)

func main() {
	logger := zap.Must(zap.NewProduction())
	config.Read("./config/config.yaml", logger)
	app.Run(logger)
}
