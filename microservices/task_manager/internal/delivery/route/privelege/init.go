package privelege

import (
	"github.com/cantylv/authorization-service/client"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/delivery/route/privelege/agent"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/delivery/route/privelege/group"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/delivery/route/privelege/privelege"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/delivery/route/privelege/user"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func InitHTTPHandlers(r *mux.Router, privelegeClient *client.Client, logger *zap.Logger) {
	agent.InitHandlers(r, privelegeClient, logger)
	user.InitHandlers(r, privelegeClient, logger)
	group.InitHandlers(r, privelegeClient, logger)
	privelege.InitHandlers(r, privelegeClient, logger)
}
