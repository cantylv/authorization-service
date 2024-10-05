package route

import (
	"net/http"

	"github.com/cantylv/authorization-service/client"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/delivery/route/archive"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/delivery/route/ping"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/delivery/route/privelege"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/middlewares"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func InitHTTPHandlers(r *mux.Router, privelegeClient *client.Client, logger *zap.Logger) http.Handler {
	s := r.PathPrefix("/api/v1").Subrouter()
	ping.InitHandler(r)
	privelege.InitHTTPHandlers(s, privelegeClient, logger)
	archive.InitHTTPHandlers(s, archiveClient, logger)
	return middlewares.Init(s, logger)
}
