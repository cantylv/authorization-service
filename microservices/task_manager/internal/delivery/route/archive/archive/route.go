package archive

import (
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/clients"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/delivery/archive/archive"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func InitHandlers(r *mux.Router, cluster *clients.Cluster, logger *zap.Logger) {
	proxyManager := archive.NewArchiveProxyManager(logger, cluster)
	r.HandleFunc("/archive/who_asks/{email_ask}", proxyManager.GetArchive).Methods("GET")
}
