package archive

import (
	"github.com/cantylv/authorization-service/microservices/archive_manager/internal/delivery/archive"
	rArchive "github.com/cantylv/authorization-service/microservices/archive_manager/internal/repo/archive"
	uArchive "github.com/cantylv/authorization-service/microservices/archive_manager/internal/usecase/archive"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func InitHandlers(r *mux.Router, clientPostgres *pgx.Conn, logger *zap.Logger) {
	repoArchive := rArchive.NewRepoLayer(clientPostgres)
	usecaseArchive := uArchive.NewUsecaseLayer(repoArchive)
	archiveManager := archive.NewHandlerArchiveManager(logger, usecaseArchive)
	r.HandleFunc("/archives", archiveManager.GetArchive).Methods("GET")
}
