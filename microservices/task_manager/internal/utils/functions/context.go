package functions

import (
	"net/http"

	pClient "github.com/cantylv/authorization-service/client"
	aClient "github.com/cantylv/authorization-service/microservices/archive_manager/client"
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

func GetCtxRequestMeta(r *http.Request) (pClient.RequestMeta, error) {
	meta, ok := r.Context().Value(mc.AccessKey(mc.RequestMeta)).(pClient.RequestMeta)
	if !ok {
		return pClient.RequestMeta{
			RealIp: uuid.NewV4().String(), // we need to specify real ip, because microservice 'privelege' uses it for log id in bad cases
		}, me.ErrNoMetaInContext
	}
	return meta, nil
}

func GetCtxRequestMetaForArchive(r *http.Request) (aClient.RequestMeta, error) {
	meta, ok := r.Context().Value(mc.AccessKey(mc.RequestMeta)).(aClient.RequestMeta)
	if !ok {
		return aClient.RequestMeta{
			RealIp: uuid.NewV4().String(), // we need to specify real ip, because microservice 'privelege' uses it for log id in bad cases
		}, me.ErrNoMetaInContext
	}
	return meta, nil
}
