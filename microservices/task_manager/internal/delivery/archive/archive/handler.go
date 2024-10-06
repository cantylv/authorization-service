package archive

import (
	"net/http"

	pClient "github.com/cantylv/authorization-service/client"
	aClient "github.com/cantylv/authorization-service/microservices/archive_manager/client"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/clients"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/entity/dto"
	f "github.com/cantylv/authorization-service/microservices/task_manager/internal/utils/functions"
	mc "github.com/cantylv/authorization-service/microservices/task_manager/internal/utils/myconstants"
	me "github.com/cantylv/authorization-service/microservices/task_manager/internal/utils/myerrors"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ArchiveProxyManager struct {
	logger          *zap.Logger
	archiveClient   *aClient.Client
	privelegeClient *pClient.Client
}

// NewArchiveProxyManager возвращает прокси менеджер, отвечающий за проксирование запросов к агентам.
func NewArchiveProxyManager(logger *zap.Logger, cluster *clients.Cluster) *ArchiveProxyManager {
	return &ArchiveProxyManager{
		logger:          logger,
		archiveClient:   cluster.ArchiveClient,
		privelegeClient: cluster.PrivelegeClient,
	}
}

func (h *ArchiveProxyManager) GetArchive(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	metaPrivelege, err := f.GetCtxRequestMeta(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	metaArchive, err := f.GetCtxRequestMetaForArchive(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	pathVars := mux.Vars(r)
	emailAsk := pathVars["email_ask"]
	// убедимся, что пользователь имеет доступ к архиву
	canExecute, status := h.privelegeClient.Privelege.CanUserExecute(emailAsk, "archive", &metaPrivelege)
	if status.Err != nil {
		h.logger.Info(status.Err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: status.Err.Error()}, status.StatusCode)
		return
	}
	if !canExecute {
		h.logger.Info(me.ErrUserDoesntHaveEnoughPrivelege.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrUserDoesntHaveEnoughPrivelege.Error()}, http.StatusForbidden)
		return
	}
	agent, reqStatus := h.archiveClient.GetArchive(&metaArchive)
	if reqStatus.Err != nil {
		h.logger.Info(reqStatus.Err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: reqStatus.Err.Error()}, reqStatus.StatusCode)
		return
	}
	f.Response(w, agent, reqStatus.StatusCode)
}
