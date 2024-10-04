package manager

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func InitHandlers(r *mux.Router, logger *zap.Logger) {
	r.HandleFunc("/api/v1/load_balancing", func(w http.ResponseWriter, r *http.Request) {})
}
