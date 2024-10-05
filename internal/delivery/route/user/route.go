package user

import (
	"net/http"

	"github.com/cantylv/authorization-service/internal/delivery/user"
	rGroup "github.com/cantylv/authorization-service/internal/repo/group"
	rUser "github.com/cantylv/authorization-service/internal/repo/user"
	uUser "github.com/cantylv/authorization-service/internal/usecase/user"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// InitHandlers инициализирует обработчики запросов для работы с пользователями (получение, удаление, создание).
func InitHandlers(r *mux.Router, postgresClient *pgx.Conn, logger *zap.Logger) {
	repoUser := rUser.NewRepoLayer(postgresClient)
	repoGroup := rGroup.NewRepoLayer(postgresClient)
	ucUser := uUser.NewUsecaseLayer(repoUser, repoGroup)
	userHandlerManager := user.NewUserHandlerManager(ucUser, logger)
	// ручки, отвечающие за создание, получение и удаление пользователя
	r.HandleFunc("/users", userHandlerManager.Create).Methods("POST")                                      // создание пользователя
	r.HandleFunc("/users/{email}", userHandlerManager.Read).Methods("GET")                                 // чтение данных пользователя
	r.HandleFunc("/users/{email}/who_deletes/{email_delete}", userHandlerManager.Delete).Methods("DELETE") // удаление пользователя
	r.HandleFunc("/openid/callback", func(http.ResponseWriter, *http.Request) {}).Methods("POST")          // callback URL для openID провайдера
}
