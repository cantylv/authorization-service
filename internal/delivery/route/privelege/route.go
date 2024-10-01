package privelege

import (
	dPrivelege "github.com/cantylv/authorization-service/internal/delivery/privelege"
	rAgent "github.com/cantylv/authorization-service/internal/repo/agent"
	rGroup "github.com/cantylv/authorization-service/internal/repo/group"
	rPrivelege "github.com/cantylv/authorization-service/internal/repo/privelege"
	rUser "github.com/cantylv/authorization-service/internal/repo/user"
	uPrivelege "github.com/cantylv/authorization-service/internal/usecase/privelege"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// InitHandlers инициализирует обработчики запросов, отвечающих за права пользователя к ресурсу.
// Пользователи принадлежат группам, в свою очередь права присваиваются группам, поэтому пользователь, находящийся
// в какой-то группе наследует ее права.
func InitHandlers(r *mux.Router, postgresClient *pgx.Conn, logger *zap.Logger) {
	repoAgent := rAgent.NewRepoLayer(postgresClient)
	repoPrivelege := rPrivelege.NewRepoLayer(postgresClient)
	repoUser := rUser.NewRepoLayer(postgresClient)
	repoGroup := rGroup.NewRepoLayer(postgresClient)
	usecasePrivelege := uPrivelege.NewUsecaseLayer(repoAgent, repoPrivelege, repoUser, repoGroup)
	privelegeHandlerManager := dPrivelege.NewPrivelegeHandlerManager(usecasePrivelege, logger)
	r.HandleFunc("/api/v1/groups/{group_name}/priveleges/new/agents/{agent_name}/who_adds/{email_add}", privelegeHandlerManager.AddAgentToGroup).Methods("POST")                 // добавляет группе нового агента
	r.HandleFunc("/api/v1/groups/{group_name}/priveleges/delete/agents/{agent_name}/who_deletes/{email_delete}", privelegeHandlerManager.DeleteAgentFromGroup).Methods("DELETE") // удаляет у группы агента
	r.HandleFunc("/api/v1/groups/{group_name}/priveleges/who_asks/{email_ask}", privelegeHandlerManager.GetGroupAgents).Methods("GET")                                           // возвращает список агентов группы
	r.HandleFunc("/api/v1/users/{email}/check_access/agents/{agent_name}", privelegeHandlerManager.CanUserExecute).Methods("GET")                                                // проверяет, можно ли пользователю пользоваться агентом
}
