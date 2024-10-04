package route

import (
	"net/http"

	"github.com/cantylv/authorization-service/microservices/task_manager/internal/delivery/route/manager"
	"github.com/cantylv/authorization-service/microservices/task_manager/internal/middlewares"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func InitHttpHandlers(r *mux.Router, logger *zap.Logger) http.Handler {
	manager.InitHandlers(r, logger)
	return middlewares.Init(r, logger)
}
