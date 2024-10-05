package archive

import (
	"context"
	"database/sql"
	"errors"

	ent "github.com/cantylv/authorization-service/microservices/archive_manager/internal/entity"
	"github.com/cantylv/authorization-service/microservices/archive_manager/internal/repo/archive"
	"github.com/cantylv/authorization-service/microservices/archive_manager/internal/utils/myerrors"
)

type Usecase interface {
	GetArchive(ctx context.Context) ([]*ent.Record, error)
}

var _ Usecase = (*UsecaseLayer)(nil)

type UsecaseLayer struct {
	repoArchive archive.Repo
}

func NewUsecaseLayer(repoArchive archive.Repo) *UsecaseLayer {
	return &UsecaseLayer{
		repoArchive: repoArchive,
	}
}

func (u *UsecaseLayer) GetArchive(ctx context.Context) ([]*ent.Record, error) {
	records, err := u.repoArchive.Get(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myerrors.ErrNoArchive
		}
		return nil, err
	}
	return records, nil
}
