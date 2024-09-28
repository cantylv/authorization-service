package group

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// InitHandlers инициализирует обработчики запросов, отвечающих за права пользователя к ресурсу
func InitHandlers(r *mux.Router, postgresClient *pgx.Conn, logger *zap.Logger) {
	r.HandleFunc("/api/v1/groups/{groupID}/add_user/{email}/who_invites/{email_invite}", func(http.ResponseWriter, *http.Request) {}).Methods("POST") // добавляет пользователя в группу
	r.HandleFunc("/api/v1/users/{email}/groups/who_asks/{email_ask}", func(http.ResponseWriter, *http.Request) {}).Methods("GET")                     // возвращает список групп пользователя
	r.HandleFunc("/api/v1/groups/{groupID}/kick_user/{email}/who_kicks/{email_kicks}", func(http.ResponseWriter, *http.Request) {}).Methods("PUT")    // удаляет пользователя из группы
	r.HandleFunc("/api/v1/users/{email}/groups/new", func(http.ResponseWriter, *http.Request) {}).Methods("POST")                                     // добавляет заявку на создание группы
	r.HandleFunc("/api/v1/users/{email}/groups/{groupID}/status", func(http.ResponseWriter, *http.Request) {}).Methods("POST")                        // подтверждает/отклоняет заявку на создание группы
}
