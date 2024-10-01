package agent

import (
	"errors"
	"net/http"

	"github.com/asaskevich/govalidator"
	ent "github.com/cantylv/authorization-service/internal/entity"
	"github.com/cantylv/authorization-service/internal/entity/dto"
	"github.com/cantylv/authorization-service/internal/usecase/agent"
	f "github.com/cantylv/authorization-service/internal/utils/functions"
	mc "github.com/cantylv/authorization-service/internal/utils/myconstants"
	me "github.com/cantylv/authorization-service/internal/utils/myerrors"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type AgentHandlerManager struct {
	logger       *zap.Logger
	usecaseAgent agent.Usecase
}

// NewAgentHandlerManager возвращает менеджер хендлеров, отвечающих за создание агентов. Работают только для root
// пользователя.
func NewAgentHandlerManager(usecaseAgent agent.Usecase, logger *zap.Logger) *AgentHandlerManager {
	return &AgentHandlerManager{
		logger:       logger,
		usecaseAgent: usecaseAgent,
	}
}

// AddAgent добавляет агента, который обрабатывает сооответствующие ему запросы
// Создать агента может только root
func (h *AgentHandlerManager) CreateAgent(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}

	pathVars := mux.Vars(r)
	agentName := pathVars["agent_name"]
	emailCreate := pathVars["email_create"]
	if !govalidator.IsEmail(emailCreate) {
		h.logger.Info(me.ErrInvalidEmail.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidEmail.Error()}, http.StatusBadRequest)
		return
	}
	a, err := h.usecaseAgent.CreateAgent(r.Context(), emailCreate, agentName)
	if err != nil {
		if errors.Is(err, me.ErrOnlyRootCanAddAgent) {
			h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusForbidden)
			return
		}
		if errors.Is(err, me.ErrAgentAlreadyExist) {
			h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusBadRequest)
			return
		}
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}
	f.Response(w, a, http.StatusOK)
}

// DeleteAgent удаляет агента, который обрабатывает сооответствующие ему запросы
// Удалить агента может только root
func (h *AgentHandlerManager) DeleteAgent(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}

	pathVars := mux.Vars(r)
	agentName := pathVars["agent_name"]
	emailDelete := pathVars["email_delete"]
	if !govalidator.IsEmail(emailDelete) {
		h.logger.Info(me.ErrInvalidEmail.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidEmail.Error()}, http.StatusBadRequest)
		return
	}
	err = h.usecaseAgent.DeleteAgent(r.Context(), emailDelete, agentName)
	if err != nil {
		if errors.Is(err, me.ErrOnlyRootCanDeleteAgent) {
			h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusForbidden)
			return
		}
		if errors.Is(err, me.ErrAgentNotExist) {
			h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusBadRequest)
			return
		}
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}
	f.Response(w, dto.ResponseDetail{Detail: "agent was succesful deleted"}, http.StatusOK)
}

// GetAgents возвращает список всех агентов
// Получить список может только root
func (h *AgentHandlerManager) GetAgents(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}

	pathVars := mux.Vars(r)
	emailRead := pathVars["email_read"]
	if !govalidator.IsEmail(emailRead) {
		h.logger.Info(me.ErrInvalidEmail.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidEmail.Error()}, http.StatusBadRequest)
		return
	}
	a, err := h.usecaseAgent.GetAgents(r.Context(), emailRead)
	if err != nil {
		if errors.Is(err, me.ErrOnlyRootCanGetAgents) {
			h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusForbidden)
			return
		}
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}
	if a == nil {
		a = make([]*ent.Agent, 0)
	}
	f.Response(w, a, http.StatusOK)
}
