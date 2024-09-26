package middlewares

import (
	"fmt"
	"net/http"

	"github.com/cantylv/authorization-service/internal/entity/dto"
	f "github.com/cantylv/authorization-service/internal/utils/functions"
	e "github.com/cantylv/authorization-service/internal/utils/myerrors"
	"go.uber.org/zap"
)

func Recover(h http.Handler, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error(fmt.Sprintf("error while handling request: %v", err))
				f.Response(w, dto.ResponseError{Error: e.ErrInternal.Error()}, http.StatusInternalServerError)
				return
			}
		}()
		h.ServeHTTP(w, r)
	})
}
