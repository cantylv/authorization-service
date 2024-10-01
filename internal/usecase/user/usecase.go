package user

import (
	"context"
	"database/sql"
	"errors"

	ent "github.com/cantylv/authorization-service/internal/entity"
	"github.com/cantylv/authorization-service/internal/entity/dto"
	"github.com/cantylv/authorization-service/internal/repo/group"
	"github.com/cantylv/authorization-service/internal/repo/user"
	f "github.com/cantylv/authorization-service/internal/utils/functions"
	me "github.com/cantylv/authorization-service/internal/utils/myerrors"
	"github.com/spf13/viper"
)

type Usecase interface {
	Create(ctx context.Context, authData *dto.CreateData) (*ent.User, error)
	Read(ctx context.Context, email string) (*ent.User, error)
	Delete(ctx context.Context, userEmail, userEmailAsk string) error
}

var _ Usecase = (*UsecaseLayer)(nil)

type UsecaseLayer struct {
	repoUser  user.Repo
	repoGroup group.Repo
}

// NewUsecaseLayer возращает структуру уровня usecase для работы с пользователями
func NewUsecaseLayer(repoUser user.Repo, repoGroup group.Repo) *UsecaseLayer {
	return &UsecaseLayer{
		repoUser:  repoUser,
		repoGroup: repoGroup,
	}
}

// Create создает пользователя. Пароль, передаваемый в теле запроса, хэшируется с помощью соли 
// алгоритмом Argon2.
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

// Read возвращает данные о пользователе.
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

// Delete удаляет пользователя из системы.
// Нельзя удалить root пользователя, а также любого ответственного за группу. Также удалить пользователя 
// может только root, либо пользователь сам себя удаляет.
func (u *UsecaseLayer) Delete(ctx context.Context, userEmail, userEmailAsk string) error {
	if userEmail == viper.GetString("root_email") {
		return me.ErrCantDeleteRoot
	}
	// проверка существования пользователя, которого удаляем
	uDB, err := u.repoUser.GetByEmail(ctx, userEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return me.ErrUserNotExist
		}
		return err
	}
	// проверяем, что пользователь не является ответственным за организации
	groups, err := u.repoGroup.GetUserGroups(ctx, uDB.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if len(groups) != 0 {
		return me.ErrUserIsResponsible
	}
	// случай, когда пользователь удаляет сам себя
	if userEmail == userEmailAsk {
		return u.repoUser.DeleteByEmail(ctx, userEmail)
	}
	// удалить пользователя из системы может только root пользователь
	if userEmailAsk != viper.GetString("root_email") {
		return me.ErrOnlyRootCanDeleteUser
	}
	return u.repoUser.DeleteByEmail(ctx, userEmail)
}
