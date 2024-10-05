package history

import (
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/delivery/privelege/agent"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func InitHandlers(r *mux.Router, logger *zap.Logger) {
	proxyManager := agent.NewAgentProxyManager(logger, archiveClient)
	r.HandleFunc("/archives/who_asks/{email_ask}", proxyManager.CreateAgent).Methods("GET")
}
