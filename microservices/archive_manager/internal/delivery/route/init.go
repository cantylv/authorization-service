package route

import (
	"net/http"

	"github.com/cantylv/authorization-service/microservices/archive_manager/internal/delivery/route/archive"
	"github.com/cantylv/authorization-service/microservices/archive_manager/internal/delivery/route/ping"
	"github.com/cantylv/authorization-service/microservices/archive_manager/internal/middlewares"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func InitHTTPHandlers(r *mux.Router, postgresClient *pgx.Conn, logger *zap.Logger) http.Handler {
	s := r.PathPrefix("/api/v1").Subrouter()
	ping.InitHandlers(s)
	archive.InitHandlers(s, postgresClient, logger)
	return middlewares.Init(s, logger)
}
