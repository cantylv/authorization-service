package user

import (
	"net/http"

	"github.com/cantylv/authorization-service/internal/delivery/user"
	rUser "github.com/cantylv/authorization-service/internal/repo/user"
	uUser "github.com/cantylv/authorization-service/internal/usecase/user"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// InitHandlers инициализирует обработчики запросов для работы с пользователями (получение, удаление, создание).
func InitHandlers(r *mux.Router, postgresClient *pgx.Conn, logger *zap.Logger) {
	repoUser := rUser.NewRepoLayer(postgresClient)
	ucUser := uUser.NewUsecaseLayer(repoUser)
	userHandlerManager := user.NewUserHandlerManager(ucUser, logger)
	// ручки, отвечающие за создание, получение и удаление пользователя
	r.HandleFunc("/api/v1/users", userHandlerManager.Create).Methods("POST")                                       // создание пользователя
	r.HandleFunc("/api/v1/users/{email}/who_asks/{email_ask}", userHandlerManager.Read).Methods("GET")             // чтение данных пользователя
	r.HandleFunc("/api/v1/users/{email}/who_deletes/{email_deleted}", userHandlerManager.Delete).Methods("DELETE") // удаление пользователя
	r.HandleFunc("/api/v1/openid/callback", func(http.ResponseWriter, *http.Request) {}).Methods("POST")           // callback URL для openID провайдера
}
