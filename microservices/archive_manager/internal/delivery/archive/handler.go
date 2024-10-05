package archive

import (
	"errors"
	"net/http"

	"github.com/cantylv/authorization-service/microservices/archive_manager/internal/entity/dto"
	"github.com/cantylv/authorization-service/microservices/archive_manager/internal/usecase/archive"
	f "github.com/cantylv/authorization-service/microservices/archive_manager/internal/utils/functions"
	mc "github.com/cantylv/authorization-service/microservices/archive_manager/internal/utils/myconstants"
	"github.com/cantylv/authorization-service/microservices/archive_manager/internal/utils/myerrors"
	"go.uber.org/zap"
)

type HandlerArchiveManager struct {
	logger    *zap.Logger
	ucArchive archive.Usecase
}

func NewHandlerArchiveManager(logger *zap.Logger, ucArchive archive.Usecase) *HandlerArchiveManager {
	return &HandlerArchiveManager{
		logger:    logger,
		ucArchive: ucArchive,
	}
}

func (h *HandlerArchiveManager) GetArchive(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	records, err := h.ucArchive.GetArchive(r.Context())
	if err != nil {
		if errors.Is(err, myerrors.ErrNoArchive) {
			h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusBadRequest)
			return
		}
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: myerrors.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}
	f.Response(w, records, http.StatusOK)
}
