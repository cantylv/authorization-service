package role

import (
	// "github.com/cantylv/authorization-service/internal/usecase/role"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/cantylv/authorization-service/internal/entity/dto"
	"github.com/cantylv/authorization-service/internal/usecase/role"
	"github.com/cantylv/authorization-service/internal/usecase/user"
	f "github.com/cantylv/authorization-service/internal/utils/functions"
	mc "github.com/cantylv/authorization-service/internal/utils/myconstants"
	me "github.com/cantylv/authorization-service/internal/utils/myerrors"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type RoleHandlerManager struct {
	ucUser user.Usecase
	ucRole role.Usecase
	logger *zap.Logger
}

// NewRoleHandlerManager возвращает менеджер хендлеров, отвечающих за получение прав пользователя на ресурс
func NewRoleHandlerManager(ucUser user.Usecase, ucRole role.Usecase, logger *zap.Logger) *RoleHandlerManager {
	return &RoleHandlerManager{
		ucUser: ucUser,
		logger: logger,
		ucRole: ucRole,
	}
}

func (h *RoleHandlerManager) CanUserExecute(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	pathVars := mux.Vars(r)
	processName := pathVars["process_name"]
	userEmail := pathVars["email"]
	userAskEmail := pathVars["email_ask"]
	if !govalidator.IsEmail(userEmail) || !govalidator.IsEmail(userAskEmail) {
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidEmail.Error()}, http.StatusBadRequest)
		return
	}
	isCanExecute, err := h.ucRole.CanExecute(r.Context(), userEmail, processName, userAskEmail)
	if err != nil {
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}
	if !isCanExecute {

	}
	f.Response(w, dto.ResponseDetail{Detail: "user has rules to run the task"}, http.StatusOK)
}
