package user

import (
	"net/http"

	"github.com/cantylv/authorization-service/client"
	f "github.com/cantylv/authorization-service/internal/utils/functions"
	mc "github.com/cantylv/authorization-service/internal/utils/myconstants"
	"go.uber.org/zap"
)

type UserProxyManager struct {
	logger          *zap.Logger
	privelegeClient *client.Client
}

// NewUserProxyManager возвращает прокси менеджер, отвечающий за создание/удаление пользователя из системы
func NewUserProxyManager(logger *zap.Logger, privelegeClient *client.Client) *UserProxyManager {
	return &UserProxyManager{
		logger: logger,
		privelegeClient: privelegeClient,
	}
}

func (h *UserProxyManager) Create(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
}

func (h *UserProxyManager) Read(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
}

func (h *UserProxyManager) Delete(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
}
