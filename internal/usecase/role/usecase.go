package role

import (
	"context"
	"database/sql"
	"errors"

	"github.com/cantylv/authorization-service/internal/repo/role"
	"github.com/cantylv/authorization-service/internal/repo/user"
	me "github.com/cantylv/authorization-service/internal/utils/myerrors"
)

type Usecase interface {
	CanExecute(ctx context.Context, userEmail, processName, userAskEmail string) (bool, error)
}

var _ Usecase = (*UsecaseLayer)(nil)

type UsecaseLayer struct {
	repoRole role.Repo
	repoUser user.Repo
}

func NewUsecaseLayer(repoRole role.Repo, repoUser user.Repo) *UsecaseLayer {
	return &UsecaseLayer{
		repoRole: repoRole,
		repoUser: repoUser,
	}
}

func (u *UsecaseLayer) CanExecute(ctx context.Context, userEmail, processName, userAskEmail string) (bool, error) {
	// проверяем, существует ли пользователь, права которого хотим проверить
	_, err := u.repoUser.GetByEmail(ctx, userEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, me.ErrUserNotExist
		}
		return false, err
	}
	// проверяем, существует ли пользователь, который запрашивает права
	_, err = u.repoUser.GetByEmail(ctx, userEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, me.ErrUserNotExist
		}
		return false, err
	}
	// нужно проверить, является ли uAskEmail тем, кто может дать грант
	// пользователь может дать грант, если он является владельцем группы
	// и группа владеет микросервисомц
	// дать доступ к микросервису может только root
	return true, nil
}
