package agent

import (
	"github.com/cantylv/authorization-service/client"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/delivery/agent"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func InitHandlers(r *mux.Router, privelegeClient *client.Client, logger *zap.Logger) {
	proxyManager := agent.NewAgentProxyManager(logger, privelegeClient)
	r.HandleFunc("/agents/{agent_name}/who_creates/{email_create}", proxyManager.CreateAgent).Methods("POST")
	r.HandleFunc("/agents/{agent_name}/who_deletes/{email_delete}", proxyManager.DeleteAgent).Methods("DELETE")
	r.HandleFunc("/agents/who_reads/{email_read}", proxyManager.GetAgents).Methods("GET")
}
