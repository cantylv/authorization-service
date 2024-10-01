package route

import (
	"net/http"

	"github.com/cantylv/authorization-service/internal/delivery/route/agent"
	"github.com/cantylv/authorization-service/internal/delivery/route/group"
	"github.com/cantylv/authorization-service/internal/delivery/route/user"
	"github.com/cantylv/authorization-service/internal/middlewares"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// InitHTTPHandlers инициализирует обработчики запросов, а также добавляет цепочку middlewares в обработку запроса.
func InitHTTPHandlers(r *mux.Router, postgresClient *pgx.Conn, logger *zap.Logger) http.Handler {
	user.InitHandlers(r, postgresClient, logger)
	group.InitHandlers(r, postgresClient, logger)
	agent.InitHandlers(r, postgresClient, logger)
	h := middlewares.Init(r, logger)
	return h
}
