package archive

import (
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/clients"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/delivery/route/archive/archive"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func InitHTTPHandlers(r *mux.Router, cluster *clients.Cluster, logger *zap.Logger) {
	archive.InitHandlers(r, cluster, logger)
}
