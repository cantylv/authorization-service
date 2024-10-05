package route

import (
	"net/http"

	"github.com/cantylv/authorization-service/client"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/delivery/route/privelege/agent"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/delivery/route/privelege/group"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/delivery/route/privelege/ping"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/delivery/route/privelege/privelege"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/delivery/route/privelege/user"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/middlewares"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func InitHTTPHandlers(r *mux.Router, privelegeClient *client.Client, logger *zap.Logger) http.Handler {
	s := r.PathPrefix("/api/v1").Subrouter()
	ping.InitHandlers(s)
	agent.InitHandlers(s, privelegeClient, logger)
	user.InitHandlers(s, privelegeClient, logger)
	group.InitHandlers(s, privelegeClient, logger)
	privelege.InitHandlers(s, privelegeClient, logger)
	return middlewares.Init(s, logger)
}
