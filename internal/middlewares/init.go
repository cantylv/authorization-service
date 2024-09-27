package middlewares

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func Init(r *mux.Router, logger *zap.Logger) (h http.Handler) {
	h = Cors(r)
	h = Recover(h, logger)
	h = Access(h, logger)
	return h
}
