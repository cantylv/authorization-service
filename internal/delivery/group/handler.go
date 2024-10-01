package group

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	ent "github.com/cantylv/authorization-service/internal/entity"
	"github.com/cantylv/authorization-service/internal/entity/dto"
	"github.com/cantylv/authorization-service/internal/usecase/group"
	f "github.com/cantylv/authorization-service/internal/utils/functions"
	mc "github.com/cantylv/authorization-service/internal/utils/myconstants"
	me "github.com/cantylv/authorization-service/internal/utils/myerrors"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type GroupHandlerManager struct {
	logger       *zap.Logger
	usecaseGroup group.Usecase
}

// NewGroupHandlerManager возвращает менеджер хендлеров, отвечающих за создание групп пользователей, добавление
// пользователей в группы и удаление из них. Заявка на создание группы и принятие/отклонение её root пользователем.
func NewGroupHandlerManager(usecaseGroup group.Usecase, logger *zap.Logger) *GroupHandlerManager {
	return &GroupHandlerManager{
		logger:       logger,
		usecaseGroup: usecaseGroup,
	}
}

// AddUserToGroup позволяет добавить пользователя в группу. Добавить в группу может только ответственный за группу пользователь.
func (h *GroupHandlerManager) AddUserToGroup(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	pathVars := mux.Vars(r)
	groupName := pathVars["group_name"]
	userEmail := pathVars["email"]
	if !govalidator.IsEmail(userEmail) {
		h.logger.Info(me.ErrInvalidEmail.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidEmail.Error()}, http.StatusBadRequest)
		return
	}
	inviteUserEmail := pathVars["email_invite"]
	if !govalidator.IsEmail(inviteUserEmail) {
		h.logger.Info(me.ErrInvalidEmail.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidEmail.Error()}, http.StatusBadRequest)
		return
	}

	groupName, err = h.usecaseGroup.AddUserToGroup(r.Context(), userEmail, inviteUserEmail, groupName)
	if err != nil {
		if errors.Is(err, me.ErrGroupNotExist) ||
			errors.Is(err, me.ErrUserEmailMustBeDiff) ||
			errors.Is(err, me.ErrUserNotExist) ||
			errors.Is(err, me.ErrOnlyOwnerCanAddUserToGroup) ||
			errors.Is(err, me.ErrUserAlreadyInGroup) {
			h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusBadRequest)
			return
		}
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}

	f.Response(w, dto.ResponseDetail{Detail: fmt.Sprintf("user was succesful added to group '%s'", groupName)}, http.StatusOK)
}

// GetUserGroups возвращает группы пользователя. Получить группы пользователя может получить любой пользователь,
// но ему покажутся только общие группы. То есть если пользователь А имеет группы users, devs, а пользователь B - users,
// то при запросе групп пользователя А пользователем В ему отобразится только одна группа users.
func (h *GroupHandlerManager) GetUserGroups(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	pathVars := mux.Vars(r)
	userEmail := pathVars["email"]
	if !govalidator.IsEmail(userEmail) {
		h.logger.Info(me.ErrInvalidEmail.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidEmail.Error()}, http.StatusBadRequest)
		return
	}
	askUserEmail := pathVars["email_ask"]
	if !govalidator.IsEmail(askUserEmail) {
		h.logger.Info(me.ErrInvalidEmail.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidEmail.Error()}, http.StatusBadRequest)
		return
	}
	groups, err := h.usecaseGroup.GetUserGroups(r.Context(), userEmail, askUserEmail)
	if err != nil {
		if errors.Is(err, me.ErrUserNotExist) {
			h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusBadRequest)
			return
		}
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}
	if groups == nil {
		groups = make([]*ent.Group, 0)
	}
	f.Response(w, groups, http.StatusOK)
}

func (h *GroupHandlerManager) KickOutUser(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	pathVars := mux.Vars(r)
	groupName := pathVars["group_name"]
	userEmail := pathVars["email"]
	if !govalidator.IsEmail(userEmail) {
		h.logger.Info(me.ErrInvalidEmail.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidEmail.Error()}, http.StatusBadRequest)
		return
	}
	kickUserEmail := pathVars["email_kick"]
	if !govalidator.IsEmail(kickUserEmail) {
		h.logger.Info(me.ErrInvalidEmail.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidEmail.Error()}, http.StatusBadRequest)
		return
	}
	groupName, err = h.usecaseGroup.KickUserFromGroup(r.Context(), userEmail, kickUserEmail, groupName)
	if err != nil {
		if errors.Is(err, me.ErrUserNotExist) ||
			errors.Is(err, me.ErrGroupNotExist) ||
			errors.Is(err, me.ErrUserIsNotInGroup) ||
			errors.Is(err, me.ErrOnlyOwnerCanDeleteUserFromGroup) ||
			errors.Is(err, me.ErrOwnerCantExitFromGroup) {
			h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusBadRequest)
			return
		}
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}
	f.Response(w, dto.ResponseDetail{Detail: fmt.Sprintf("user was succesful deleted from group '%s'", groupName)}, http.StatusOK)
}

// RequestToCreateGroup создает заявку пользователя на создание группы пользователей, у которой будут
// в будущем свои права на выполнение различных процессов.
func (h *GroupHandlerManager) RequestToCreateGroup(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	pathVars := mux.Vars(r)
	groupName := pathVars["group_name"]
	userEmail := pathVars["email_add"]
	if !govalidator.IsEmail(userEmail) {
		h.logger.Info(me.ErrInvalidEmail.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidEmail.Error()}, http.StatusBadRequest)
		return
	}
	bid, err := h.usecaseGroup.MakeRequestToCreateGroup(r.Context(), userEmail, groupName)
	if err != nil {
		if errors.Is(err, me.ErrUserNotExist) ||
			errors.Is(err, me.ErrGroupAlreadyExist) ||
			errors.Is(err, me.ErrBidAlreadyExist) {
			h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusBadRequest)
			return
		}
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}

	f.Response(w, bid, http.StatusOK)
}

func (h *GroupHandlerManager) ChangeBidStatus(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	pathVars := mux.Vars(r)
	groupName := pathVars["group_name"]
	userEmail := pathVars["email"]
	if !govalidator.IsEmail(userEmail) {
		h.logger.Info(me.ErrInvalidEmail.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidEmail.Error()}, http.StatusBadRequest)
		return
	}
	userChangeStatus := pathVars["email_change_status"]
	if !govalidator.IsEmail(userChangeStatus) {
		h.logger.Info(me.ErrInvalidEmail.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidEmail.Error()}, http.StatusBadRequest)
		return
	}
	bidStatus := r.URL.Query().Get("status")
	bid, err := h.usecaseGroup.UpdateRequestStatus(r.Context(), userEmail, groupName, userChangeStatus, bidStatus)
	if err != nil {
		if errors.Is(err, me.ErrInvalidStatus) ||
			errors.Is(err, me.ErrOnlyRootCanChangeBidStatus) ||
			errors.Is(err, me.ErrUserNotExist) ||
			errors.Is(err, me.ErrBidNotExist) ||
			errors.Is(err, me.ErrGroupAlreadyExist) {
			h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusBadRequest)
			return
		}
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}
	f.Response(w, bid, http.StatusOK)
}

func (h *GroupHandlerManager) ChangeOwner(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	pathVars := mux.Vars(r)
	groupName := pathVars["group_name"]
	userEmail := pathVars["email"]
	if !govalidator.IsEmail(userEmail) {
		h.logger.Info(me.ErrInvalidEmail.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidEmail.Error()}, http.StatusBadRequest)
		return
	}
	userChangeOwner := pathVars["email_change_owner"]
	if !govalidator.IsEmail(userChangeOwner) {
		h.logger.Info(me.ErrInvalidEmail.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidEmail.Error()}, http.StatusBadRequest)
		return
	}
	g, err := h.usecaseGroup.ChangeOwner(r.Context(), userEmail, groupName, userChangeOwner)
	if err != nil {
		if errors.Is(err, me.ErrGroupNotExist) ||
			errors.Is(err, me.ErrUserNotExist) ||
			errors.Is(err, me.ErrUserIsAlreadyOwner) {
			h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusBadRequest)
			return
		}
		if errors.Is(err, me.ErrOnlyOwnerCanAppointNewOwner) {
			h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusForbidden)
			return
		}
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}
	f.Response(w, g, http.StatusOK)
}
