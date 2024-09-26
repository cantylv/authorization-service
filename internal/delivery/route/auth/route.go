package auth

import (
	dAuth "github.com/cantylv/authorization-service/internal/delivery/auth"
	rSession "github.com/cantylv/authorization-service/internal/repo/session"
	rUser "github.com/cantylv/authorization-service/internal/repo/user"
	ucAuth "github.com/cantylv/authorization-service/internal/usecase/auth"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func InitHandlers(r *mux.Router, postgresClient *pgx.Conn, logger *zap.Logger) {
	repoSession := rSession.NewRepoLayer(postgresClient)
	repoUser := rUser.NewRepoLayer(postgresClient)
	ucAuth := ucAuth.NewUsecaseLayer(repoSession, repoUser)
	authHandlerManager := dAuth.NewAuthHandlerManager(ucAuth, logger)
	r.HandleFunc("/api/v1/auth/signin", authHandlerManager.SignIn).Methods("POST")
	r.HandleFunc("/api/v1/auth/signup", authHandlerManager.SignUp).Methods("POST")
	r.HandleFunc("/api/v1/auth/signout", authHandlerManager.SignOut).Methods("POST")
}
