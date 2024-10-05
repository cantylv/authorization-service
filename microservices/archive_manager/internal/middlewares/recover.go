package middlewares

import (
	"fmt"
	"net/http"

	"github.com/cantylv/authorization-service/microservices/archive_manager/internal/entity/dto"
	f "github.com/cantylv/authorization-service/microservices/archive_manager/internal/utils/functions"
	me "github.com/cantylv/authorization-service/microservices/archive_manager/internal/utils/myerrors"
	"go.uber.org/zap"
)

// Recover middleware для обработки паники, возникающей в работе сервера. В случае паники возвращается
// json-объект c сообщением об ошибке внутри сервера и статусом 500.
func Recover(h http.Handler, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error(fmt.Sprintf("error while handling request: %v", err))
				f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
				return
			}
		}()
		h.ServeHTTP(w, r)
	})
}
