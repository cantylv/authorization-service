package manager

import (
	"fmt"
	"net/http"

	f "github.com/cantylv/authorization-service/microservices/task_manager/internal/utils/functions"
	mc "github.com/cantylv/authorization-service/microservices/task_manager/internal/utils/myconstants"
	"github.com/gorilla/mux"
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
	// request processing:
	// получаем запрос
	// относим его к конкретному агенту
	// делаем запрос к микросервису проверки прав
	// получаем ответ
	// делаем запрос к целевому ресурсу, если можно

	// у нас есть всего один дефолтный агент privelege. root пользователь должен добавить агенты.
	// предположим, что у нас пока что только 2 микросервиса: тот, что отвечает за пользователя и права к агентам
	// и archive-manager. Для доступа к последнему делается 'preflight' запрос к первому микросервису, чтобы
	// узнать, имеет ли пользователь доступ к архиву.

	// '/api/v1/archive' - ручка микросервиса архивов
	// для доступа к микросервису проверки прав не нужно делать предварительный запрос
	reqPath := mux.Vars(r)["urlPath"]
	
}
