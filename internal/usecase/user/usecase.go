package user

import (
	"context"
	"database/sql"
	"errors"

	ent "github.com/cantylv/authorization-service/internal/entity"
	"github.com/cantylv/authorization-service/internal/entity/dto"
	"github.com/cantylv/authorization-service/internal/repo/user"
	f "github.com/cantylv/authorization-service/internal/utils/functions"
	me "github.com/cantylv/authorization-service/internal/utils/myerrors"
)

type Usecase interface {
	Create(ctx context.Context, authData *dto.CreateData) (*ent.User, error)
	Read(ctx context.Context, email string) (*ent.User, error)
	Delete(ctx context.Context, email string) error
}

var _ Usecase = (*UsecaseLayer)(nil)

type UsecaseLayer struct {
	repoUser user.Repo
}

func NewUsecaseLayer(repoUser user.Repo) *UsecaseLayer {
	return &UsecaseLayer{
		repoUser: repoUser,
	}
}

func (u *UsecaseLayer) Create(ctx context.Context, authData *dto.CreateData) (*ent.User, error) {
	// проверяем, существует ли уже пользователь c такой почтой
	// если да, то возвращаем ошибку
	uDB, err := u.repoUser.GetByEmail(ctx, authData.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if uDB != nil {
		return nil, me.ErrUserAlreadyExist
	}
	// получаем хэшированный пароль вместе с солью
	hashedPassword, err := f.GetHashedPassword(authData.Password)
	if err != nil {
		return nil, err
	}
	userNew, err := u.repoUser.Create(ctx, newUserFromSignUpForm(authData, hashedPassword))
	if err != nil {
		return nil, err
	}
	return userNew, nil
}

func (u *UsecaseLayer) Read(ctx context.Context, email string) (*ent.User, error) {
	uDB, err := u.repoUser.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, me.ErrUserNotExist
		}
		return nil, err
	}
	return uDB, nil
}

func (u *UsecaseLayer) Delete(ctx context.Context, email string) error {
	// проверяем, существует ли пользователь
	_, err := u.repoUser.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return me.ErrUserNotExist
		}
		return err
	}
	return u.repoUser.DeleteByEmail(ctx, email)
}
