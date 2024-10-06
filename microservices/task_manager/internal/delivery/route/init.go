package route

import (
	"net/http"

	"github.com/cantylv/authorization-service/microservices/task_manager/internal/clients"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/delivery/route/archive"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/delivery/route/privelege"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/middlewares"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func InitHTTPHandlers(r *mux.Router, cluster *clients.Cluster, logger *zap.Logger) http.Handler {
	s := r.PathPrefix("/api/v1").Subrouter()
	privelege.InitHTTPHandlers(s, cluster.PrivelegeClient, logger)
	archive.InitHTTPHandlers(s, cluster, logger)
	return middlewares.Init(s, logger)
}
