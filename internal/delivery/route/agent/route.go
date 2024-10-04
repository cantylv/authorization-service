package agent

import (
	"github.com/cantylv/authorization-service/internal/delivery/agent"
	rAgent "github.com/cantylv/authorization-service/internal/repo/agent"
	ucAgent "github.com/cantylv/authorization-service/internal/usecase/agent"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// InitHandlers инициализирует обработчики запросов, отвечающих crd agent
func InitHandlers(r *mux.Router, postgresClient *pgx.Conn, logger *zap.Logger) {
	repoAgent := rAgent.NewRepoLayer(postgresClient)
	usecaseAgent := ucAgent.NewUsecaseLayer(repoAgent)
	agentHandlerManager := agent.NewAgentHandlerManager(usecaseAgent, logger)
	r.HandleFunc("/agents/{agent_name}/who_creates/{email_create}", agentHandlerManager.CreateAgent).Methods("POST")   // создает агента
	r.HandleFunc("/agents/{agent_name}/who_deletes/{email_delete}", agentHandlerManager.DeleteAgent).Methods("DELETE") // удаляет агента
	r.HandleFunc("/agents/who_reads/{email_read}", agentHandlerManager.GetAgents).Methods("GET")                       // возвращает список доступных агентов
}
