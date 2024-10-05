package agent

import (
	"net/http"

	"github.com/cantylv/authorization-service/client"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/entity/dto"
	f "github.com/cantylv/authorization-service/microservices/task_manager/internal/utils/functions"
	mc "github.com/cantylv/authorization-service/microservices/task_manager/internal/utils/myconstants"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type AgentProxyManager struct {
	logger          *zap.Logger
	privelegeClient *client.Client
}

// NewAgentProxyManager возвращает прокси менеджер, отвечающий за проксирование запросов к агентам.
func NewAgentProxyManager(logger *zap.Logger, privelegeClient *client.Client) *AgentProxyManager {
	return &AgentProxyManager{
		logger:          logger,
		privelegeClient: privelegeClient,
	}
}

func (h *AgentProxyManager) CreateAgent(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	meta, err := f.GetCtxRequestMeta(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	pathVars := mux.Vars(r)
	agentName := pathVars["agent_name"]
	emailCreate := pathVars["email_create"]
	agent, reqStatus := h.privelegeClient.Agent.Create(agentName, emailCreate, &meta)
	if reqStatus.Err != nil {
		h.logger.Info(reqStatus.Err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: reqStatus.Err.Error()}, reqStatus.StatusCode)
		return
	}
	f.Response(w, agent, reqStatus.StatusCode)
}

func (h *AgentProxyManager) DeleteAgent(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	meta, err := f.GetCtxRequestMeta(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	pathVars := mux.Vars(r)
	agentName := pathVars["agent_name"]
	emailDelete := pathVars["email_delete"]
	agent, reqStatus := h.privelegeClient.Agent.Delete(agentName, emailDelete, &meta)
	if reqStatus.Err != nil {
		h.logger.Info(reqStatus.Err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: reqStatus.Err.Error()}, reqStatus.StatusCode)
		return
	}
	f.Response(w, agent, reqStatus.StatusCode)
}

func (h *AgentProxyManager) GetAgents(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	meta, err := f.GetCtxRequestMeta(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	pathVars := mux.Vars(r)
	emailRead := pathVars["email_read"]
	agent, reqStatus := h.privelegeClient.Agent.GetAll(emailRead, &meta)
	if reqStatus.Err != nil {
		h.logger.Info(reqStatus.Err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: reqStatus.Err.Error()}, reqStatus.StatusCode)
		return
	}
	f.Response(w, agent, reqStatus.StatusCode)
}
