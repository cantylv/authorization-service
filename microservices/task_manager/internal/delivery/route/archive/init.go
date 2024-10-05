package archive

import (
	"github.com/cantylv/authorization-service/client"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/delivery/route/archive/history"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func InitHTTPHandlers(r *mux.Router, privelegeClient *client.Client, logger *zap.Logger) {
	history.InitHandlers(r, logger)
}
