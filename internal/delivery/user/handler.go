package user

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/cantylv/authorization-service/internal/entity/dto"
	"github.com/cantylv/authorization-service/internal/usecase/user"
	f "github.com/cantylv/authorization-service/internal/utils/functions"
	mc "github.com/cantylv/authorization-service/internal/utils/myconstants"
	me "github.com/cantylv/authorization-service/internal/utils/myerrors"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type UserHandlerManager struct {
	ucUser user.Usecase
	logger *zap.Logger
}

// NewUserHandlerManager возвращает менеджер хендлеров, отвечающих за создание/удаление пользователя из системы
func NewUserHandlerManager(ucUser user.Usecase, logger *zap.Logger) *UserHandlerManager {
	return &UserHandlerManager{
		ucUser: ucUser,
		logger: logger,
	}
}

// Create метод создания пользователя, в случае успеха возвращает пользователю его данные.
// Не требует идентификации в запросе, так как инициируется неавторизованным пользователем.
func (h *UserHandlerManager) Create(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidData.Error()}, http.StatusBadRequest)
		return
	}
	var signForm dto.CreateData
	err = json.Unmarshal(body, &signForm)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidData.Error()}, http.StatusBadRequest)
		return
	}
	isValid, err := f.Validate(signForm)
	if err != nil || !isValid {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidData.Error()}, http.StatusBadRequest)
		return
	}

	u, err := h.ucUser.Create(r.Context(), &signForm)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		if errors.Is(err, me.ErrUserAlreadyExist) {
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusBadRequest)
			return
		}
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}
	f.Response(w, getUserWithoutPassword(u), http.StatusOK)
}

// Read метод чтения данных пользователя, в случае успеха возвращает пользователю его данные.
// Требует идентификации в запросе, так как инициируется авторизованным пользователем.
func (h *UserHandlerManager) Read(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	userEmail := mux.Vars(r)["email"]
	if !govalidator.IsEmail(userEmail) {
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidEmail.Error()}, http.StatusBadRequest)
		return
	}
	u, err := h.ucUser.Read(r.Context(), userEmail)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		if errors.Is(err, me.ErrUserNotExist) {
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusBadRequest)
			return
		}
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}
	f.Response(w, getUserWithoutPassword(u), http.StatusOK)
}

// Read метод чтения данных пользователя, в случае успеха возвращает пользователю его данные.
// Требует идентификации в запросе, так как инициируется авторизованным пользователем.
func (h *UserHandlerManager) Delete(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	userEmail := mux.Vars(r)["email"]
	if !govalidator.IsEmail(userEmail) {
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidEmail.Error()}, http.StatusBadRequest)
		return
	}
	err = h.ucUser.Delete(r.Context(), userEmail)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		if errors.Is(err, me.ErrUserNotExist) {
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusBadRequest)
			return
		}
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}
	f.Response(w, dto.ResponseDetail{Detail: "user was succesful deleted"}, http.StatusOK)
}
