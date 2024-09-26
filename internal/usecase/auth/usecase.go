package auth

import (
	"context"
	"database/sql"
	"errors"

	ent "github.com/cantylv/authorization-service/internal/entity"
	"github.com/cantylv/authorization-service/internal/entity/dto"
	"github.com/cantylv/authorization-service/internal/repo/session"
	"github.com/cantylv/authorization-service/internal/repo/user"
	f "github.com/cantylv/authorization-service/internal/utils/functions"
	me "github.com/cantylv/authorization-service/internal/utils/myerrors"
)

type Usecase interface {
	SignIn(ctx context.Context, authData *dto.SignInData, meta *ent.Session) (*ent.User, *ent.Session, error)
	SignUp(ctx context.Context, authData *dto.SignUpData, meta *ent.Session) (*ent.User, *ent.Session, error)
	SignOut(ctx context.Context, meta *ent.Session) error
}

var _ Usecase = (*UsecaseLayer)(nil)

type UsecaseLayer struct {
	repoSession session.Repo
	repoUser    user.Repo
}

func NewUsecaseLayer(repoSession session.Repo, repoUser user.Repo) *UsecaseLayer {
	return &UsecaseLayer{
		repoSession: repoSession,
		repoUser:    repoUser,
	}
}

const maxUserSessions = 5

func (u *UsecaseLayer) SignIn(ctx context.Context, authData *dto.SignInData, meta *ent.Session) (*ent.User, *ent.Session, error) {
	// проверяем, существует ли пользователь
	uDB, err := u.repoUser.GetByEmail(ctx, authData.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, me.ErrUserNotExist
		}
		return nil, nil, err
	}
	// проверяем пароль
	if !f.IsPasswordsEqual(authData.Password, uDB.Password) {
		return nil, nil, me.ErrPasswordMismatch
	}
	// проверяем количество текущих сессий этого пользователя
	// если их больше 5, то мы их удаляем и выдаем новую сессию (подозрительно, если у пользователя много сессий)
	sessions, err := u.repoSession.ReadSessions(ctx, uDB.Id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, nil, err
	}
	if len(sessions) >= maxUserSessions {
		err = u.repoSession.DeleteUserSessions(ctx, uDB.Id)
		if err != nil {
			return nil, nil, err
		}
	}
	s, err := u.repoSession.CreateSession(ctx, uDB.Id, meta)
	if err != nil {
		return nil, nil, err
	}
	return uDB, s, nil
}

func (u *UsecaseLayer) SignUp(ctx context.Context, authData *dto.SignUpData, meta *ent.Session) (*ent.User, *ent.Session, error) {
	// проверяем, существует ли пользователь c такой почтой
	// если да, то возвращаем ошибку
	uDB, err := u.repoUser.GetByEmail(ctx, authData.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, nil, err
	}
	if uDB != nil {
		return nil, nil, me.ErrUserAlreadyExist
	}
	// получаем хэшированный пароль вместе с солью
	hashedPassword, err := f.GetHashedPassword(authData.Password)
	if err != nil {
		return nil, nil, err
	}
	userNew, err := u.repoUser.Create(ctx, newUserFromSignUpForm(authData, hashedPassword))
	if err != nil {
		return nil, nil, err
	}
	s, err := u.repoSession.CreateSession(ctx, userNew.Id, meta)
	if err != nil {
		return nil, nil, err
	}
	return userNew, s, nil
}

func (u *UsecaseLayer) SignOut(ctx context.Context, meta *ent.Session) error {
	err := u.repoSession.DeleteSessionByToken(ctx, meta.RefreshToken)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}
