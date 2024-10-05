package ping

import (
	"net/http"

	"github.com/cantylv/authorization-service/microservices/task_manager/internal/entity/dto"
	f "github.com/cantylv/authorization-service/microservices/task_manager/internal/utils/functions"

	"github.com/gorilla/mux"
)

func InitHandler(r *mux.Router) {
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		f.Response(w, dto.ResponseDetail{Detail: "pong"}, http.StatusOK)
	}).Methods("GET")
}
