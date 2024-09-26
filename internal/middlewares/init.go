package middlewares

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func Init(r *mux.Router, logger *zap.Logger) (h http.Handler) {
	h = JwtVerification(r, logger)
	// h = Csrf(h)
	h = Cors(h)
	h = Recover(h, logger)
	h = Access(h, logger)
	return h
}
