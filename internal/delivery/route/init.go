package route

import (
	"net/http"

	"github.com/cantylv/authorization-service/internal/delivery/route/agent"
	"github.com/cantylv/authorization-service/internal/delivery/route/group"
	"github.com/cantylv/authorization-service/internal/delivery/route/ping"
	"github.com/cantylv/authorization-service/internal/delivery/route/privelege"
	"github.com/cantylv/authorization-service/internal/delivery/route/user"
	"github.com/cantylv/authorization-service/internal/middlewares"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// InitHTTPHandlers инициализирует обработчики запросов, а также добавляет цепочку middlewares в обработку запроса.
func InitHTTPHandlers(r *mux.Router, postgresClient *pgx.Conn, logger *zap.Logger) http.Handler {
	s := r.PathPrefix("/api/v1").Subrouter()
	ping.InitHandlers(s)
	agent.InitHandlers(s, postgresClient, logger)
	user.InitHandlers(s, postgresClient, logger)
	group.InitHandlers(s, postgresClient, logger)
	privelege.InitHandlers(s, postgresClient, logger)
	return middlewares.Init(s, logger)
}
