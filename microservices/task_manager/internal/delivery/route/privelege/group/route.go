package group

import (
	"github.com/cantylv/authorization-service/client"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/delivery/group"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func InitHandlers(r *mux.Router, privelegeClient *client.Client, logger *zap.Logger) {
	proxyManager := group.NewGroupProxyManager(logger, privelegeClient)
	r.HandleFunc("/groups/{group_name}/add_user/{email}/who_invites/{email_invite}", proxyManager.AddUserToGroup).Methods("POST")
	r.HandleFunc("/users/{email}/groups/who_asks/{email_ask}", proxyManager.GetUserGroups).Methods("GET")
	r.HandleFunc("/groups/{group_name}/kick_user/{email}/who_kicks/{email_kick}", proxyManager.KickOutUser).Methods("POST")
	r.HandleFunc("/groups/{group_name}/who_adds/{email_add}", proxyManager.RequestToCreateGroup).Methods("POST")
	r.HandleFunc("/users/{email}/groups/{group_name}/who_change_status/{email_change_status}", proxyManager.ChangeBidStatus).Methods("PUT")
	r.HandleFunc("/groups/{group_name}/users/{email}/who_change_owner/{email_change_owner}", proxyManager.ChangeOwner).Methods("PUT")
}
