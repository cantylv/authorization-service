package group

import (
	dGroup "github.com/cantylv/authorization-service/internal/delivery/group"
	repoGroup "github.com/cantylv/authorization-service/internal/repo/group"
	repoUser "github.com/cantylv/authorization-service/internal/repo/user"
	"github.com/cantylv/authorization-service/internal/usecase/group"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// InitHandlers инициализирует обработчики запросов, отвечающих за права пользователя к ресурсу
func InitHandlers(r *mux.Router, postgresClient *pgx.Conn, logger *zap.Logger) {
	repoUser := repoUser.NewRepoLayer(postgresClient)
	repoGroup := repoGroup.NewRepoLayer(postgresClient)
	usecaseGroup := group.NewUsecaseLayer(repoUser, repoGroup)
	userHandlerManager := dGroup.NewGroupHandlerManager(usecaseGroup, logger)
	r.HandleFunc("/groups/{group_name}/add_user/{email}/who_invites/{email_invite}", userHandlerManager.AddUserToGroup).Methods("POST")           // добавляет пользователя в группу
	r.HandleFunc("/users/{email}/groups/who_asks/{email_ask}", userHandlerManager.GetUserGroups).Methods("GET")                                   // возвращает список групп пользователя
	r.HandleFunc("/groups/{group_name}/kick_user/{email}/who_kicks/{email_kick}", userHandlerManager.KickOutUser).Methods("POST")                 // удаляет пользователя из группы
	r.HandleFunc("/groups/{group_name}/who_adds/{email_add}", userHandlerManager.RequestToCreateGroup).Methods("POST")                            // добавляет заявку на создание группы
	r.HandleFunc("/users/{email}/groups/{group_name}/who_change_status/{email_change_status}", userHandlerManager.ChangeBidStatus).Methods("PUT") // подтверждает/отклоняет заявку на создание группы ? доступна только root
	r.HandleFunc("/groups/{group_name}/users/{email}/who_change_owner/{email_change_owner}", userHandlerManager.ChangeOwner).Methods("PUT")       // изменяет ответственного за группу
}
