package privelege

import (
	"github.com/cantylv/authorization-service/client"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/delivery/privelege/privelege"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func InitHandlers(r *mux.Router, privelegeClient *client.Client, logger *zap.Logger) {
	proxyManager := privelege.NewPrivelegeProxyManager(logger, privelegeClient)
	// привелегии, которые назначаются группам
	r.HandleFunc("/groups/{group_name}/priveleges/new/agents/{agent_name}/who_adds/{email_add}", proxyManager.AddAgentToGroup).Methods("POST")
	r.HandleFunc("/groups/{group_name}/priveleges/delete/agents/{agent_name}/who_deletes/{email_delete}", proxyManager.DeleteAgentFromGroup).Methods("DELETE")
	r.HandleFunc("/groups/{group_name}/priveleges/who_asks/{email_ask}", proxyManager.GetGroupAgents).Methods("GET")
	// привелегии, которые назначаются конкретному пользователю
	r.HandleFunc("/users/{email}/priveleges/new/agents/{agent_name}/who_adds/{email_add}", proxyManager.AddAgentToUser).Methods("POST")
	r.HandleFunc("/users/{email}/priveleges/delete/agents/{agent_name}/who_deletes/{email_delete}", proxyManager.DeleteAgentFromUser).Methods("DELETE")
	r.HandleFunc("/users/{email}/priveleges/who_asks/{email_ask}", proxyManager.GetUserAgents).Methods("GET")
	r.HandleFunc("/users/{email}/check_access/agents/{agent_name}", proxyManager.CanUserExecute).Methods("GET")
}
