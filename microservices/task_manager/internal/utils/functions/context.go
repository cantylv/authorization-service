package functions

import (
	"net/http"

	"github.com/cantylv/authorization-service/client"
	mc "github.com/cantylv/authorization-service/microservices/task_manager/internal/utils/myconstants"
	me "github.com/cantylv/authorization-service/microservices/task_manager/internal/utils/myerrors"
	"github.com/satori/uuid"
)

func GetCtxRequestID(r *http.Request) (string, error) {
	requestID, ok := r.Context().Value(mc.AccessKey(mc.RequestID)).(string)
	if !ok {
		// we need to authenticate requests using unique keys | remote address is OK
		return r.RemoteAddr, me.ErrNoRequestIdInContext
	}
	return requestID, nil
}

func GetCtxRequestMeta(r *http.Request) (client.RequestMeta, error) {
	meta, ok := r.Context().Value(mc.AccessKey(mc.RequestMeta)).(client.RequestMeta)
	if !ok {
		return client.RequestMeta{
			RealIp: uuid.NewV4().String(), // we need to specify real ip, because microservice 'privelege' uses it for log id in bad cases
		}, me.ErrNoMetaInContext
	}
	return meta, nil
}
