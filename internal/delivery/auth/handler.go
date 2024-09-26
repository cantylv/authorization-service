package tokens

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/cantylv/authorization-service/internal/entity/dto"
	ucAuth "github.com/cantylv/authorization-service/internal/usecase/auth"
	f "github.com/cantylv/authorization-service/internal/utils/functions"
	mc "github.com/cantylv/authorization-service/internal/utils/myconstants"
	me "github.com/cantylv/authorization-service/internal/utils/myerrors"
	"go.uber.org/zap"
)

type AuthHandlerManager struct {
	ucAuth ucAuth.Usecase
	logger *zap.Logger
}

// NewAuthHandlerManager возвращает менеджер хендлеров авторизации/аутентификации
func NewAuthHandlerManager(ucAuth ucAuth.Usecase, logger *zap.Logger) *AuthHandlerManager {
	return &AuthHandlerManager{
		ucAuth: ucAuth,
		logger: logger,
	}
}

// SignIn метод авторизации пользователя, в случае успеха возвращает пользователю его данные и jwt_token в теле ответа.
// Проставляет куку refresh_token на домен авторизации сервера (/api/v1/auth).
// Если идет запрос с jwt_token-ом, то запрос возвращает ошибку со статусом
func (h *AuthHandlerManager) SignIn(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	// проверяем наличие jwt_token-а
	// если он есть в контексте, значит он валидный --> пользователь авторизован --> такой запрос мы отклоняем
	_, ok := r.Context().Value(mc.AccessKey(mc.JwtPayload)).(*dto.JwtPayload)
	if ok {
		h.logger.Error(me.ErrAlreadyAuthorized.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrAlreadyAuthorized.Error()}, http.StatusForbidden)
		return
	}
	meta, err := getMetadataFromConnection(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		if errors.Is(err, me.ErrInvalidRemoteIp) || errors.Is(err, me.ErrInvalidUserAgent) {
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusBadRequest)
			return
		}
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidData.Error()}, http.StatusBadRequest)
		return
	}
	var signInForm dto.SignInForm
	err = json.Unmarshal(body, &signInForm)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidData.Error()}, http.StatusBadRequest)
		return
	}
	isValid, err := f.Validate(signInForm.Payload)
	if err != nil || !isValid {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidData.Error()}, http.StatusBadRequest)
		return
	}

	u, s, err := h.ucAuth.SignIn(r.Context(), &signInForm.Payload, meta)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}
	f.SetCookie(w, r, s)
	jwtToken, err := createJwtToken(s)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}
	f.Response(w, getUserResponse(jwtToken, u), http.StatusOK)
}

// SignUp метод регистрации пользователя, в случае успеха возвращает пользователю его данные и jwt_token в теле ответа.
// Проставляет куку refresh_token на домен авторизации сервера (/api/v1/auth).
// Если идет запрос с jwt_token-ом, то запрос возвращает ошибку со статусом.
func (h *AuthHandlerManager) SignUp(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	// проверяем наличие jwt_token-а
	// если он есть в контексте, значит он валидный --> пользователь авторизован --> мы принимаем запрос
	// в любом другом случае отклоняем запрос
	_, ok := r.Context().Value(mc.AccessKey(mc.JwtPayload)).(*dto.JwtPayload)
	if !ok {
		h.logger.Error(me.ErrAlreadyAuthorized.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrAlreadyAuthorized.Error()}, http.StatusForbidden)
		return
	}
	meta, err := getMetadataFromConnection(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		if errors.Is(err, me.ErrInvalidRemoteIp) || errors.Is(err, me.ErrInvalidUserAgent) {
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusBadRequest)
			return
		}
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}
	if meta.RefreshToken == "" {
		h.logger.Error(me.ErrNoRefreshToken.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrNoRefreshToken.Error()}, http.StatusForbidden)
		return
	}

	err = h.ucAuth.SignOut(r.Context(), meta)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}
	f.CookieExpired(w, r)
	f.Response(w, dto.ResponseDetail{Detail: "you were succesful signed out"}, http.StatusOK)
}

func (h *AuthHandlerManager) SignOut(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	isExist, err := f.IsCookieExist(r, mc.RefreshToken)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}
	if !isExist {
		h.logger.Error(me.ErrNoRefreshToken.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrNoRefreshToken.Error()}, http.StatusUnauthorized)
		return
	}
	meta, err := getMetadataFromConnection(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		if errors.Is(err, me.ErrInvalidRemoteIp) || errors.Is(err, me.ErrInvalidUserAgent) {
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusBadRequest)
			return
		}
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidData.Error()}, http.StatusBadRequest)
		return
	}
	var signForm *dto.SignUpData
	err = json.Unmarshal(body, signForm)
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

	u, s, err := h.ucAuth.SignUp(r.Context(), signForm, meta)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}
	f.SetCookie(w, r, s)
	jwtToken, err := createJwtToken(s)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}
	f.Response(w, getUserResponse(jwtToken, u), http.StatusOK)
}
