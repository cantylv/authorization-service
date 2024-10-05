package user

import (
	"net/http"

	"github.com/cantylv/authorization-service/client"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/delivery/user"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func InitHandlers(r *mux.Router, privelegeClient *client.Client, logger *zap.Logger) {
	proxyManager := user.NewUserProxyManager(logger, privelegeClient)
	r.HandleFunc("/users", proxyManager.Create).Methods("POST")
	r.HandleFunc("/users/{email}", proxyManager.Read).Methods("GET")
	r.HandleFunc("/users/{email}/who_deletes/{email_ask}", proxyManager.Delete).Methods("DELETE")
	r.HandleFunc("/openid/callback", func(http.ResponseWriter, *http.Request) {}).Methods("POST")
}
