package ping

import (
	"net/http"

	"github.com/cantylv/authorization-service/microservices/archive_manager/internal/entity/dto"
	f "github.com/cantylv/authorization-service/microservices/archive_manager/internal/utils/functions"
	"github.com/gorilla/mux"
)

func InitHandlers(r *mux.Router) {
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		f.Response(w, dto.ResponseDetail{Detail: "pong"}, http.StatusOK)
	}).Methods("GET")
}
