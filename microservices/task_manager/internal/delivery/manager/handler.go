package manager

import (
	"net/http"

	f "github.com/cantylv/authorization-service/microservices/task_manager/internal/utils/functions"
	mc "github.com/cantylv/authorization-service/microservices/task_manager/internal/utils/myconstants"
	"go.uber.org/zap"
)

type ManagerHttpRequestsHadler struct {
	logger *zap.Logger
}

func NewManagerHttpRequestsHadler(logger *zap.Logger) *ManagerHttpRequestsHadler {
	return &ManagerHttpRequestsHadler{
		logger: logger,
	}
}

func (h *ManagerHttpRequestsHadler) Load(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	// получаем запрос
	// относим его к конкретному агенту
	// делаем запрос к микросервису проверки прав
	// получаем ответ
	// делаем запрос к целевому ресурсу, если можно

	// у нас есть всего один дефолтный агент privelege
}
