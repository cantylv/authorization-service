package privelege

import (
	// "github.com/cantylv/authorization-service/internal/usecase/role"

	"net/http"

	"github.com/cantylv/authorization-service/client"
	f "github.com/cantylv/authorization-service/internal/utils/functions"
	mc "github.com/cantylv/authorization-service/internal/utils/myconstants"
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
}

func (h *PrivelegeProxyManager) DeleteAgentFromGroup(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
}

func (h *PrivelegeProxyManager) GetGroupAgents(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
}

func (h *PrivelegeProxyManager) AddAgentToUser(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
}

func (h *PrivelegeProxyManager) DeleteAgentFromUser(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
}

func (h *PrivelegeProxyManager) GetUserAgents(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
}

func (h *PrivelegeProxyManager) CanUserExecute(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
}
