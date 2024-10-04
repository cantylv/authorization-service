package ping

import (
	"net/http"

	"github.com/cantylv/authorization-service/internal/entity/dto"
	"github.com/cantylv/authorization-service/internal/utils/functions"
	"github.com/gorilla/mux"
)

func InitHandlers(r *mux.Router) {
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		functions.Response(w, dto.ResponseDetail{Detail: "pong"}, http.StatusOK)
	}).Methods("GET")
}
