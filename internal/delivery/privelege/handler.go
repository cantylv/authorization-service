package privelege

import (
	// "github.com/cantylv/authorization-service/internal/usecase/role"
	"errors"
	"net/http"

	"github.com/asaskevich/govalidator"
	ent "github.com/cantylv/authorization-service/internal/entity"
	"github.com/cantylv/authorization-service/internal/entity/dto"
	"github.com/cantylv/authorization-service/internal/usecase/privelege"
	f "github.com/cantylv/authorization-service/internal/utils/functions"
	mc "github.com/cantylv/authorization-service/internal/utils/myconstants"
	me "github.com/cantylv/authorization-service/internal/utils/myerrors"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type PrivelegeHandlerManager struct {
	ucPrivelege privelege.Usecase
	logger      *zap.Logger
}

// NewPrivelegeHandlerManager возвращает менеджер хендлеров, отвечающих за получение прав пользователя на ресурс.
// Права на запуск агента добавляются группе, а пользователь уже наследует от нее права.
func NewPrivelegeHandlerManager(ucPrivelege privelege.Usecase, logger *zap.Logger) *PrivelegeHandlerManager {
	return &PrivelegeHandlerManager{
		ucPrivelege: ucPrivelege,
		logger:      logger,
	}
}

func (h *PrivelegeHandlerManager) AddAgentToGroup(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	pathVars := mux.Vars(r)
	groupName := pathVars["group_name"]
	agentName := pathVars["agent_name"]
	emailAdd := pathVars["email_add"]
	if !govalidator.IsEmail(emailAdd) {
		h.logger.Info(me.ErrInvalidEmail.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidEmail.Error()}, http.StatusBadRequest)
		return
	}
	err = h.ucPrivelege.AddAgentToGroup(r.Context(), agentName, groupName, emailAdd)
	if err != nil {
		if errors.Is(err, me.ErrAgentNotExist) ||
			errors.Is(err, me.ErrGroupNotExist) ||
			errors.Is(err, me.ErrGroupAgentAlreadyExist) {
			h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusBadRequest)
			return
		}
		if errors.Is(err, me.ErrOnlyRootCanAddAgent) {
			h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusForbidden)
			return
		}
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}

	f.Response(w, dto.ResponseDetail{Detail: "agent was succesful added to group"}, http.StatusOK)
}

func (h *PrivelegeHandlerManager) DeleteAgentFromGroup(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	pathVars := mux.Vars(r)
	groupName := pathVars["group_name"]
	agentName := pathVars["agent_name"]
	emailDelete := pathVars["email_delete"]
	if !govalidator.IsEmail(emailDelete) {
		h.logger.Info(me.ErrInvalidEmail.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidEmail.Error()}, http.StatusBadRequest)
		return
	}
	err = h.ucPrivelege.DeleteAgentFromGroup(r.Context(), agentName, groupName, emailDelete)
	if err != nil {
		if errors.Is(err, me.ErrAgentNotExist) ||
			errors.Is(err, me.ErrGroupNotExist) ||
			errors.Is(err, me.ErrGroupAgentNotExist) {
			h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusBadRequest)
			return
		}
		if errors.Is(err, me.ErrOnlyRootCanDeleteAgent) {
			h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusForbidden)
			return
		}
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}

	f.Response(w, dto.ResponseDetail{Detail: "agent was succesful deleted from group"}, http.StatusOK)
}

func (h *PrivelegeHandlerManager) GetGroupAgents(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	pathVars := mux.Vars(r)
	groupName := pathVars["group_name"]
	emailAsk := pathVars["email_ask"]
	if !govalidator.IsEmail(emailAsk) {
		h.logger.Info(me.ErrInvalidEmail.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidEmail.Error()}, http.StatusBadRequest)
		return
	}
	agents, err := h.ucPrivelege.GetGroupAgents(r.Context(), groupName, emailAsk)
	if err != nil {
		if errors.Is(err, me.ErrGroupNotExist) ||
			errors.Is(err, me.ErrUserNotExist) {
			h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusBadRequest)
			return
		}
		if errors.Is(err, me.ErrUserIsNotOwner) {
			h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusForbidden)
			return
		}
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}
	if agents == nil {
		agents = make([]*ent.Agent, 0)
	}
	f.Response(w, agents, http.StatusOK)
}

func (h *PrivelegeHandlerManager) AddAgentToUser(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	pathVars := mux.Vars(r)
	agentName := pathVars["agent_name"]
	email := pathVars["email"]
	if !govalidator.IsEmail(email) {
		h.logger.Info(me.ErrInvalidEmail.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidEmail.Error()}, http.StatusBadRequest)
		return
	}
	emailAdd := pathVars["email_add"]
	if !govalidator.IsEmail(emailAdd) {
		h.logger.Info(me.ErrInvalidEmail.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidEmail.Error()}, http.StatusBadRequest)
		return
	}
	err = h.ucPrivelege.AddAgentToGroup(r.Context(), agentName, email, emailAdd)
	if err != nil {
		if errors.Is(err, me.ErrAgentNotExist) ||
			errors.Is(err, me.ErrUserNotExist) ||
			errors.Is(err, me.ErrUserAgentAlreadyExist) {
			h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusBadRequest)
			return
		}
		if errors.Is(err, me.ErrOnlyRootCanAddAgent) {
			h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusForbidden)
			return
		}
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}

	f.Response(w, dto.ResponseDetail{Detail: "agent was succesful added to user"}, http.StatusOK)
}

func (h *PrivelegeHandlerManager) DeleteAgentFromUser(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	pathVars := mux.Vars(r)
	agentName := pathVars["agent_name"]
	email := pathVars["email"]
	if !govalidator.IsEmail(email) {
		h.logger.Info(me.ErrInvalidEmail.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidEmail.Error()}, http.StatusBadRequest)
		return
	}
	emailDelete := pathVars["email_delete"]
	if !govalidator.IsEmail(emailDelete) {
		h.logger.Info(me.ErrInvalidEmail.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidEmail.Error()}, http.StatusBadRequest)
		return
	}
	err = h.ucPrivelege.DeleteAgentFromUser(r.Context(), agentName, email, emailDelete)
	if err != nil {
		if errors.Is(err, me.ErrAgentNotExist) ||
			errors.Is(err, me.ErrUserNotExist) ||
			errors.Is(err, me.ErrUserAgentNotExist) {
			h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusBadRequest)
			return
		}
		if errors.Is(err, me.ErrOnlyRootCanDeleteAgent) {
			h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusForbidden)
			return
		}
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}

	f.Response(w, dto.ResponseDetail{Detail: "agent was succesful deleted from group"}, http.StatusOK)
}

func (h *PrivelegeHandlerManager) GetUserAgents(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	pathVars := mux.Vars(r)
	email := pathVars["email"]
	if !govalidator.IsEmail(email) {
		h.logger.Info(me.ErrInvalidEmail.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidEmail.Error()}, http.StatusBadRequest)
		return
	}
	emailAsk := pathVars["email_ask"]
	if !govalidator.IsEmail(emailAsk) {
		h.logger.Info(me.ErrInvalidEmail.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidEmail.Error()}, http.StatusBadRequest)
		return
	}
	agents, err := h.ucPrivelege.GetUserAgents(r.Context(), email, emailAsk)
	if err != nil {
		if errors.Is(err, me.ErrUserNotExist) {
			h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusBadRequest)
			return
		}
		if errors.Is(err, me.ErrGetUserAgents) {
			h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusForbidden)
			return
		}
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}
	if agents == nil {
		agents = make([]*ent.Agent, 0)
	}
	f.Response(w, agents, http.StatusOK)
}

func (h *PrivelegeHandlerManager) CanUserExecute(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	pathVars := mux.Vars(r)
	agentName := pathVars["agent_name"]
	userEmail := pathVars["email"]
	if !govalidator.IsEmail(userEmail) {
		h.logger.Info(me.ErrInvalidEmail.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidEmail.Error()}, http.StatusBadRequest)
		return
	}
	canExecute, err := h.ucPrivelege.CanExecute(r.Context(), userEmail, agentName)
	if err != nil {
		if errors.Is(err, me.ErrAgentNotExist) ||
			errors.Is(err, me.ErrUserNotExist) {
			h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusBadRequest)
			return
		}
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}
	if canExecute {
		f.Response(w, map[string]bool{"can_execute": true}, http.StatusOK)
		return
	}
	f.Response(w, map[string]bool{"can_execute": false}, http.StatusOK)
}
