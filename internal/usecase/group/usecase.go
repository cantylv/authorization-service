package group

import (
	"context"
	"database/sql"
	"errors"

	ent "github.com/cantylv/authorization-service/internal/entity"
	"github.com/cantylv/authorization-service/internal/entity/dto"
	"github.com/cantylv/authorization-service/internal/repo/group"
	"github.com/cantylv/authorization-service/internal/repo/user"
	mc "github.com/cantylv/authorization-service/internal/utils/myconstants"
	me "github.com/cantylv/authorization-service/internal/utils/myerrors"
	"github.com/spf13/viper"
)

type Usecase interface {
	AddUserToGroup(ctx context.Context, userEmail, inviteUserEmail, groupName string) (string, error)
	GetUserGroups(ctx context.Context, userEmail, askUserEmail string) ([]*ent.Group, error)
	KickUserFromGroup(ctx context.Context, userEmail, kickUserEmail, groupName string) (string, error)
	MakeRequestToCreateGroup(ctx context.Context, userEmail, groupName string) (*dto.Bid, error)
	UpdateRequestStatus(ctx context.Context, userEmail, groupName, userChangeStatus, status string) (*dto.Bid, error)
	ChangeOwner(ctx context.Context, userEmail, groupName, userChangeOwnerEmail string) (*ent.Group, error)
}

var _ Usecase = (*UsecaseLayer)(nil)

type UsecaseLayer struct {
	repoUser  user.Repo
	repoGroup group.Repo
}

// NewUsecaseLayer возвращает структуру уровня usecase, управляющую группами пользователей
func NewUsecaseLayer(repoUser user.Repo, repoGroup group.Repo) *UsecaseLayer {
	return &UsecaseLayer{
		repoUser:  repoUser,
		repoGroup: repoGroup,
	}
}

// AddUserToGroup позволяет добавить пользователя в группу. Добавить в группу может только основатель этой группы.
// Метод возвращает название группы, в которую пользоваетель был добавлен и ошибку в случае неудачи.
func (u *UsecaseLayer) AddUserToGroup(ctx context.Context, userEmail, inviteUserEmail, groupName string) (string, error) {
	// проверяем, существует ли группа, в которую мы хотим добавить пользователя
	groupDB, err := u.repoGroup.GetGroup(ctx, groupName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", me.ErrGroupNotExist
		}
		return "", err
	}
	// проверяем, разные это пользователи или нет (на дурака)
	if userEmail == inviteUserEmail {
		return "", me.ErrUserEmailMustBeDiff
	}
	// проверяем, существует ли пользователь, которого собираемся добавить в группу
	uDB, err := u.repoUser.GetByEmail(ctx, userEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", me.ErrUserNotExist
		}
		return "", err
	}
	// проверяем, есть ли пользователь, который собирается добавить в группу
	uInviter, err := u.repoUser.GetByEmail(ctx, inviteUserEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", me.ErrUserNotExist
		}
		return "", err
	}
	// проверяем, является ли пользователь, который добавляет в группу другого пользователя,
	// ответственным за нее (создателем другими словами). Также root пользователь может добавить кого угодно в любую группу.
	if inviteUserEmail != viper.GetString("root_email") && groupDB.OwnerID != uInviter.ID {
		return "", me.ErrOnlyOwnerCanAddUserToGroup
	}
	// проверка на то, есть ли уже пользователь в этой группе
	isParticipant, err := u.repoGroup.IsParticipantOfGroup(ctx, uDB.ID, groupDB.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}
	if isParticipant {
		return "", me.ErrUserAlreadyInGroup
	}
	// добавляем пользователя | добавление происходит без подтверждения root пользователя
	// в данном случае мы рассчитываем, что создателям групп можно доверять
	err = u.repoGroup.AddUserToGroup(ctx, uDB.ID, groupDB.ID)
	if err != nil {
		return "", err
	}
	return groupDB.Name, nil
}

// GetUserGroups возвращает список групп пользователя. Показывает только общие группы с другими пользователями.
func (u *UsecaseLayer) GetUserGroups(ctx context.Context, userEmail, askUserEmail string) ([]*ent.Group, error) {
	// проверяем, существует ли пользователь, чьи группы мы хотим получить
	uDB, err := u.repoUser.GetByEmail(ctx, userEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, me.ErrUserNotExist
		}
		return nil, err
	}
	// проверяем, разные это пользователи или нет
	if userEmail != askUserEmail {
		// проверяем, есть ли пользователь, который хочет получить список групп пользователя
		uInviter, err := u.repoUser.GetByEmail(ctx, askUserEmail)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, me.ErrUserNotExist
			}
			return nil, err
		}
		if askUserEmail != viper.GetString("root_email") && userEmail != viper.GetString("root_email") {
			groups, err := u.repoGroup.GetCommonGroups(ctx, uDB.ID, uInviter.ID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, nil
				}
				return nil, err
			}
			return groups, nil
		}
	}
	groups, err := u.repoGroup.GetUserGroups(ctx, uDB.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return groups, nil
}

// KickUserFromGroup удаляет пользователя из группы
func (u *UsecaseLayer) KickUserFromGroup(ctx context.Context, userEmail, kickUserEmail, groupName string) (string, error) {
	// проверяем, существует ли группа, из которую мы хотим удалить пользователя
	groupDB, err := u.repoGroup.GetGroup(ctx, groupName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", me.ErrGroupNotExist
		}
		return "", err
	}
	// root пользователь присутствует во всех группах, его нельзя от туда удалить
	if userEmail == viper.GetString("root_email") {
		return "", me.ErrDeleteRootFromGroup
	}
	// проверяем, существует ли пользователь, которого собираемся удалить из группы
	uDB, err := u.repoUser.GetByEmail(ctx, userEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", me.ErrUserNotExist
		}
		return "", err
	}
	// проверяем, есть ли пользователь в этой группе
	_, err = u.repoGroup.IsParticipantOfGroup(ctx, uDB.ID, groupDB.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", me.ErrUserIsNotInGroup
		}
		return "", err
	}
	// ограничение: владелец группы не может выйти из беседы, для того чтобы покинуть, необходимо назначить нового владельца
	if groupDB.OwnerID == uDB.ID {
		return "", me.ErrOwnerCantExitFromGroup
	}
	// проверяем, пользователь сам покидает группу или нет
	if userEmail == kickUserEmail {
		err = u.repoGroup.KickUserFromGroup(ctx, uDB.ID, groupDB.ID)
		if err != nil {
			return "", err
		}
		return groupName, nil
	}
	// проверяем, есть ли пользователь, который собирается удалить пользователя из группы
	uKicker, err := u.repoUser.GetByEmail(ctx, kickUserEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", me.ErrUserNotExist
		}
		return "", err
	}
	// пользователя из группы может удалить только владелец группы
	// проверим, что это так и есть | не забываем, что root пользователь может также удалить
	if kickUserEmail != viper.GetString("root_email") && uKicker.ID != groupDB.OwnerID {
		return "", me.ErrOnlyOwnerCanDeleteUserFromGroup
	}
	err = u.repoGroup.KickUserFromGroup(ctx, uDB.ID, groupDB.ID)
	if err != nil {
		return "", err
	}
	return groupName, nil
}

// MakeRequestToCreateGroup создает заявку на создание группы, статус заявки "in_progress"
func (u *UsecaseLayer) MakeRequestToCreateGroup(ctx context.Context, userEmail, groupName string) (*dto.Bid, error) {
	// проверяем, существует ли пользователь, который подает заявку на создание группы
	uDB, err := u.repoUser.GetByEmail(ctx, userEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, me.ErrUserNotExist
		}
		return nil, err
	}
	// проверяем, существует ли уже группа с таким именем
	groupDB, err := u.repoGroup.GetGroup(ctx, groupName)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if groupDB != nil {
		return nil, me.ErrGroupAlreadyExist
	}
	// проверяем, существует ли уже заявка с таким именем
	bidDB, err := u.repoGroup.GetBid(ctx, uDB.ID, groupName)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if bidDB != nil {
		return nil, me.ErrBidAlreadyExist
	}
	// если запрос на создание делает root, то мы не добавляем заявку в таблицу, а сразу создаем группу
	if uDB.Email == viper.GetString("root_email") {
		groupNew, err := u.repoGroup.CreateGroup(ctx, uDB.ID, groupName)
		if err != nil {
			return nil, err
		}
		return newBidFromExistingGroup(groupNew), nil
	}
	// создаем заявку
	bid, err := u.repoGroup.MakeBidGroupCreation(ctx, uDB.ID, groupName)
	if err != nil {
		return nil, err
	}
	return bid, nil
}

func (u *UsecaseLayer) UpdateRequestStatus(ctx context.Context, userEmail, groupName, userChangeStatus, status string) (*dto.Bid, error) {
	// проверим, что статус имеет допустимое значение
	if _, ok := mc.AllowedStatus[status]; !ok {
		return nil, me.ErrInvalidStatus
	}
	// проверяем, root это или нет
	if viper.GetString("root_email") != userChangeStatus {
		return nil, me.ErrOnlyRootCanChangeBidStatus
	}
	// нужно получить id пользователя, который хочет создать новую группу (быть ее ответственным)
	uDB, err := u.repoUser.GetByEmail(ctx, userEmail)
	if err != nil {
		return nil, me.ErrUserNotExist
	}
	// проверяем, существует ли заявка с такой группой
	bidDB, err := u.repoGroup.GetBid(ctx, uDB.ID, groupName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, me.ErrBidNotExist
		}
		return nil, err
	}
	// если root пользователь отказывает в создании группы, то нам нет смысла
	// узнавать, есть ли такая группа уже
	if status == "rejected" {
		b, err := u.repoGroup.RejectGroupCreation(ctx, bidDB.ID)
		if err != nil {
			return nil, err
		}
		return b, nil
	}
	// проверяем, что в существующих группах нет такого же имени
	groupDB, err := u.repoGroup.GetGroup(ctx, groupName)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if groupDB != nil {
		return nil, me.ErrGroupAlreadyExist
	}
	// нужно получить id root пользователя
	userRoot, err := u.repoUser.GetByEmail(ctx, viper.GetString("root_email"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, me.ErrUserNotExist
		}
		return nil, err
	}
	g, err := u.repoGroup.ApproveGroupCreation(ctx, bidDB.UserId, userRoot.ID, groupName)
	if err != nil {
		return nil, err
	}
	return newBidFromExistingGroup(g), nil
}

func (u *UsecaseLayer) ChangeOwner(ctx context.Context, userEmail, groupName, userChangeOwnerEmail string) (*ent.Group, error) {
	// проверим существование группы
	groupDB, err := u.repoGroup.GetGroup(ctx, groupName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, me.ErrGroupNotExist
		}
		return nil, err
	}
	// проверим существование пользователей
	userNewOwner, err := u.repoUser.GetByEmail(ctx, userEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, me.ErrUserNotExist
		}
		return nil, err
	}
	if groupDB.OwnerID == userNewOwner.ID {
		return nil, me.ErrUserIsAlreadyOwner
	}
	// root пользователь может менять ответственного любой группы
	if userChangeOwnerEmail == viper.GetString("root_email") {
		return u.repoGroup.UpdateOwner(ctx, groupDB.ID, userNewOwner.ID)
	}
	userOldOwner, err := u.repoUser.GetByEmail(ctx, userChangeOwnerEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, me.ErrUserNotExist
		}
		return nil, err
	}
	// проверим, что пользователь является ответственным за группу
	_, err = u.repoGroup.IsOwnerOfGroup(ctx, userOldOwner.ID, groupName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, me.ErrOnlyOwnerCanAppointNewOwner
		}
		return nil, err
	}
	return u.repoGroup.UpdateOwner(ctx, groupDB.ID, userNewOwner.ID)
}
