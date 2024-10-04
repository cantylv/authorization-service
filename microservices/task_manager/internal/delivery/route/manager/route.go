package manager

import (
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/delivery/manager"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func InitHandlers(r *mux.Router, logger *zap.Logger) {
	managerHandler := manager.NewManagerHttpRequestsHadler(logger)
	r.HandleFunc("/{urlPath:.*}", managerHandler.Load)
}
