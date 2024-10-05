package user

import (
	"net/http"

	"github.com/cantylv/authorization-service/client"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/entity/dto"
	f "github.com/cantylv/authorization-service/microservices/task_manager/internal/utils/functions"
	mc "github.com/cantylv/authorization-service/microservices/task_manager/internal/utils/myconstants"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type UserProxyManager struct {
	logger          *zap.Logger
	privelegeClient *client.Client
}

// NewUserProxyManager возвращает прокси менеджер, отвечающий за создание/удаление пользователя из системы
func NewUserProxyManager(logger *zap.Logger, privelegeClient *client.Client) *UserProxyManager {
	return &UserProxyManager{
		logger:          logger,
		privelegeClient: privelegeClient,
	}
}

func (h *UserProxyManager) Create(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	meta, err := f.GetCtxRequestMeta(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	user, reqStatus := h.privelegeClient.User.Create(r.Body, &meta)
	if reqStatus.Err != nil {
		h.logger.Info(reqStatus.Err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: reqStatus.Err.Error()}, reqStatus.StatusCode)
		return
	}
	f.Response(w, user, reqStatus.StatusCode)
}

func (h *UserProxyManager) Read(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	meta, err := f.GetCtxRequestMeta(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	pathVars := mux.Vars(r)
	email := pathVars["email"]
	user, reqStatus := h.privelegeClient.User.Get(email, &meta)
	if reqStatus.Err != nil {
		h.logger.Info(reqStatus.Err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: reqStatus.Err.Error()}, reqStatus.StatusCode)
		return
	}
	f.Response(w, user, reqStatus.StatusCode)
}

func (h *UserProxyManager) Delete(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	meta, err := f.GetCtxRequestMeta(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	pathVars := mux.Vars(r)
	email := pathVars["email"]
	emailDelete := pathVars["email_delete"]
	detailMsg, reqStatus := h.privelegeClient.User.Delete(email, emailDelete, &meta)
	if reqStatus.Err != nil {
		h.logger.Info(reqStatus.Err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: reqStatus.Err.Error()}, reqStatus.StatusCode)
		return
	}
	f.Response(w, detailMsg, reqStatus.StatusCode)
}
