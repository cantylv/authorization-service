package privelege

import (
	// "github.com/cantylv/authorization-service/internal/usecase/role"

	"net/http"

	"github.com/cantylv/authorization-service/client"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/entity/dto"
	f "github.com/cantylv/authorization-service/microservices/task_manager/internal/utils/functions"
	mc "github.com/cantylv/authorization-service/microservices/task_manager/internal/utils/myconstants"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type PrivelegeProxyManager struct {
	logger          *zap.Logger
	privelegeClient *client.Client
}

func NewPrivelegeProxyManager(logger *zap.Logger, privelegeClient *client.Client) *PrivelegeProxyManager {
	return &PrivelegeProxyManager{
		logger:          logger,
		privelegeClient: privelegeClient,
	}
}

func (h *PrivelegeProxyManager) AddAgentToGroup(w http.ResponseWriter, r *http.Request) {
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
	agentName := pathVars["agent_name"]
	emailAdd := pathVars["email_add"]
	detailMsg, reqStatus := h.privelegeClient.Privelege.AddAgentToGroup(groupName, agentName, emailAdd, &meta)
	if reqStatus.Err != nil {
		h.logger.Info(reqStatus.Err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: reqStatus.Err.Error()}, reqStatus.StatusCode)
		return
	}
	f.Response(w, detailMsg, reqStatus.StatusCode)
}

func (h *PrivelegeProxyManager) DeleteAgentFromGroup(w http.ResponseWriter, r *http.Request) {
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
	agentName := pathVars["agent_name"]
	emailDelete := pathVars["email_delete"]
	detailMsg, reqStatus := h.privelegeClient.Privelege.DeleteAgentFromGroup(groupName, agentName, emailDelete, &meta)
	if reqStatus.Err != nil {
		h.logger.Info(reqStatus.Err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: reqStatus.Err.Error()}, reqStatus.StatusCode)
		return
	}
	f.Response(w, detailMsg, reqStatus.StatusCode)
}

func (h *PrivelegeProxyManager) GetGroupAgents(w http.ResponseWriter, r *http.Request) {
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
	emailAsk := pathVars["email_ask"]
	agents, reqStatus := h.privelegeClient.Privelege.GetGroupAgents(groupName, emailAsk, &meta)
	if reqStatus.Err != nil {
		h.logger.Info(reqStatus.Err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: reqStatus.Err.Error()}, reqStatus.StatusCode)
		return
	}
	f.Response(w, agents, reqStatus.StatusCode)
}

func (h *PrivelegeProxyManager) AddAgentToUser(w http.ResponseWriter, r *http.Request) {
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
	agentName := pathVars["agent_name"]
	emailAdd := pathVars["email_add"]
	detailMsg, reqStatus := h.privelegeClient.Privelege.AddAgentToUser(email, agentName, emailAdd, &meta)
	if reqStatus.Err != nil {
		h.logger.Info(reqStatus.Err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: reqStatus.Err.Error()}, reqStatus.StatusCode)
		return
	}
	f.Response(w, detailMsg, reqStatus.StatusCode)
}

func (h *PrivelegeProxyManager) DeleteAgentFromUser(w http.ResponseWriter, r *http.Request) {
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
	agentName := pathVars["agent_name"]
	emailDelete := pathVars["email_delete"]
	detailMsg, reqStatus := h.privelegeClient.Privelege.DeleteAgentFromUser(email, agentName, emailDelete, &meta)
	if reqStatus.Err != nil {
		h.logger.Info(reqStatus.Err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: reqStatus.Err.Error()}, reqStatus.StatusCode)
		return
	}
	f.Response(w, detailMsg, reqStatus.StatusCode)
}

func (h *PrivelegeProxyManager) GetUserAgents(w http.ResponseWriter, r *http.Request) {
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
	agents, reqStatus := h.privelegeClient.Privelege.GetUserAgents(email, emailAsk, &meta)
	if reqStatus.Err != nil {
		h.logger.Info(reqStatus.Err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: reqStatus.Err.Error()}, reqStatus.StatusCode)
		return
	}
	f.Response(w, agents, reqStatus.StatusCode)
}

func (h *PrivelegeProxyManager) CanUserExecute(w http.ResponseWriter, r *http.Request) {
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
	agentName := pathVars["agent_name"]
	canExecute, reqStatus := h.privelegeClient.Privelege.CanUserExecute(email, agentName, &meta)
	if reqStatus.Err != nil {
		h.logger.Info(reqStatus.Err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: reqStatus.Err.Error()}, reqStatus.StatusCode)
		return
	}
	f.Response(w, map[string]bool{"can_execute": canExecute}, reqStatus.StatusCode)
}
