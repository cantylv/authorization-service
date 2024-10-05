package group

import (
	"net/http"

	"github.com/cantylv/authorization-service/client"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/entity/dto"
	f "github.com/cantylv/authorization-service/microservices/task_manager/internal/utils/functions"
	mc "github.com/cantylv/authorization-service/microservices/task_manager/internal/utils/myconstants"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type GroupProxyManager struct {
	logger          *zap.Logger
	privelegeClient *client.Client
}

func NewGroupProxyManager(logger *zap.Logger, privelegeClient *client.Client) *GroupProxyManager {
	return &GroupProxyManager{
		logger:          logger,
		privelegeClient: privelegeClient,
	}
}

func (h *GroupProxyManager) AddUserToGroup(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	meta, err := f.GetCtxRequestMeta(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	pathVars := mux.Vars(r)
	groupName := pathVars["group_name"]
	email := pathVars["email"]
	emailInvite := pathVars["email_invite"]
	detailMsg, reqStatus := h.privelegeClient.Group.AddUserToGroup(groupName, email, emailInvite, &meta)
	if reqStatus.Err != nil {
		h.logger.Info(reqStatus.Err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: reqStatus.Err.Error()}, reqStatus.StatusCode)
		return
	}
	f.Response(w, detailMsg, reqStatus.StatusCode)
}

func (h *GroupProxyManager) GetUserGroups(w http.ResponseWriter, r *http.Request) {
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
	emailAsk := pathVars["email_ask"]
	groups, reqStatus := h.privelegeClient.Group.UserList(email, emailAsk, &meta)
	if reqStatus.Err != nil {
		h.logger.Info(reqStatus.Err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: reqStatus.Err.Error()}, reqStatus.StatusCode)
		return
	}
	f.Response(w, groups, reqStatus.StatusCode)
}

func (h *GroupProxyManager) KickOutUser(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	meta, err := f.GetCtxRequestMeta(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	pathVars := mux.Vars(r)
	groupName := pathVars["group_name"]
	email := pathVars["email"]
	emailKick := pathVars["email_kick"]
	detailMsg, reqStatus := h.privelegeClient.Group.KickOutUser(groupName, email, emailKick, &meta)
	if reqStatus.Err != nil {
		h.logger.Info(reqStatus.Err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: reqStatus.Err.Error()}, reqStatus.StatusCode)
		return
	}
	f.Response(w, detailMsg, reqStatus.StatusCode)
}

func (h *GroupProxyManager) RequestToCreateGroup(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	meta, err := f.GetCtxRequestMeta(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	pathVars := mux.Vars(r)
	groupName := pathVars["group_name"]
	emailAdd := pathVars["email_add"]
	bid, reqStatus := h.privelegeClient.Group.MakeBidToCreateGroup(groupName, emailAdd, &meta)
	if reqStatus.Err != nil {
		h.logger.Info(reqStatus.Err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: reqStatus.Err.Error()}, reqStatus.StatusCode)
		return
	}
	f.Response(w, bid, reqStatus.StatusCode)
}

func (h *GroupProxyManager) ChangeBidStatus(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	meta, err := f.GetCtxRequestMeta(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	pathVars := mux.Vars(r)
	groupName := pathVars["group_name"]
	email := pathVars["email"]
	emailChangeStatus := pathVars["email_change_status"]
	newStatus := r.URL.Query().Get("status")
	bid, reqStatus := h.privelegeClient.Group.ChangeBidStatus(groupName, email, emailChangeStatus, newStatus, &meta)
	if reqStatus.Err != nil {
		h.logger.Info(reqStatus.Err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: reqStatus.Err.Error()}, reqStatus.StatusCode)
		return
	}
	f.Response(w, bid, reqStatus.StatusCode)
}

func (h *GroupProxyManager) ChangeOwner(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	meta, err := f.GetCtxRequestMeta(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	pathVars := mux.Vars(r)
	groupName := pathVars["group_name"]
	email := pathVars["email"]
	emailChangeOwner := pathVars["email_change_owner"]
	agent, reqStatus := h.privelegeClient.Group.ChangeOwner(groupName, email, emailChangeOwner, &meta)
	if reqStatus.Err != nil {
		h.logger.Info(reqStatus.Err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: reqStatus.Err.Error()}, reqStatus.StatusCode)
		return
	}
	f.Response(w, agent, reqStatus.StatusCode)
}
