package role

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// InitHandlers инициализирует обработчики запросов, отвечающих за права пользователя к ресурсу
func InitHandlers(r *mux.Router, postgresClient *pgx.Conn, logger *zap.Logger) {
	r.HandleFunc("/api/v1/users/{email}/processes/{process_name}/who_asks/{email_ask}", func(http.ResponseWriter, *http.Request) {}).Methods("GET")           // проверяет права пользователя на выполнение задачи
	r.HandleFunc("/api/v1/users/{email}/processes/who_asks/{email_ask}", func(http.ResponseWriter, *http.Request) {}).Methods("GET")                          // возвращает доступные на выполнение пользователю задачи
	r.HandleFunc("/api/v1/users/{email}/processes/{process_name}/who_allows/{email_allow}", func(http.ResponseWriter, *http.Request) {}).Methods("POST")      // добавляет право на выполнение задачи
	r.HandleFunc("/api/v1/users/{email}/processes/{process_name}/who_prohibits/{email_prohibit}", func(http.ResponseWriter, *http.Request) {}).Methods("PUT") // забирает право на выполнение задачи
}
