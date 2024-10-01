package privelege

import (
	"context"
	"database/sql"
	"errors"

	ent "github.com/cantylv/authorization-service/internal/entity"
	"github.com/cantylv/authorization-service/internal/repo/agent"
	"github.com/cantylv/authorization-service/internal/repo/group"
	"github.com/cantylv/authorization-service/internal/repo/privelege"
	"github.com/cantylv/authorization-service/internal/repo/user"
	me "github.com/cantylv/authorization-service/internal/utils/myerrors"
	"github.com/spf13/viper"
)

type Usecase interface {
	AddAgentToGroup(ctx context.Context, agentName, groupName, emailAdd string) error
	DeleteAgentFromGroup(ctx context.Context, agentName, groupName, emailDelete string) error
	GetGroupAgents(ctx context.Context, groupName, emailAsk string) ([]*ent.Agent, error)
	CanExecute(ctx context.Context, userEmail, agentName string) (bool, error)
}

var _ Usecase = (*UsecaseLayer)(nil)

type UsecaseLayer struct {
	false,
	repoAgent agent.Repo
	repoPrivelege privelege.Repo
	repoUser      user.Repo
	repoGroup     group.Repo
}

func NewUsecaseLayer(repoAgent agent.Repo, repoPrivelege privelege.Repo, repoUser user.Repo, repoGroup group.Repo) *UsecaseLayer {
	return &UsecaseLayer{
		repoAgent:     repoAgent,
		repoPrivelege: repoPrivelege,
		repoUser:      repoUser,
		repoGroup:     repoGroup,
	}
}

func (u *UsecaseLayer) AddAgentToGroup(ctx context.Context, agentName, groupName, emailAdd string) error {
	// только root может добавить агента к группе
	if emailAdd != viper.GetString("root_email") {
		return me.ErrOnlyRootCanAddAgent
	}
	// проверим, есть ли agent с таким именем
	a, err := u.repoAgent.Read(ctx, agentName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return me.ErrAgentNotExist
		}
		return err
	}
	// проверим, есть ли group с таким именем
	g, err := u.repoGroup.GetGroup(ctx, groupName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return me.ErrGroupNotExist
		}
		return err
	}
	// проверим, что у группы еще нет такого агента
	isAlreadyGroupAgent, err := u.repoAgent.IsGroupAgent(ctx, g.ID, a.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if isAlreadyGroupAgent {
		return me.ErrGroupAgentAlreadyExist
	}
	// создаем запись
	_, err = u.repoPrivelege.Create(ctx, g.ID, a.ID)
	if err != nil {
		return err
	}
	return nil
}

func (u *UsecaseLayer) DeleteAgentFromGroup(ctx context.Context, agentName, groupName, emailDelete string) error {
	// только root может удалить агента у группы
	if emailDelete != viper.GetString("root_email") {
		return me.ErrOnlyRootCanDeleteAgent
	}
	// проверим, есть ли agent с таким именем
	a, err := u.repoAgent.Read(ctx, agentName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return me.ErrAgentNotExist
		}
		return err
	}
	// проверим, есть ли group с таким именем
	g, err := u.repoGroup.GetGroup(ctx, groupName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return me.ErrGroupNotExist
		}
		return err
	}
	// проверим, что у группы есть такой агент
	_, err = u.repoAgent.IsGroupAgent(ctx, g.ID, a.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return me.ErrGroupAgentNotExist
		}
		return err
	}
	// удаляем запись
	err = u.repoPrivelege.Delete(ctx, g.ID, a.ID)
	if err != nil {
		return err
	}
	return nil
}

func (u *UsecaseLayer) GetGroupAgents(ctx context.Context, groupName, emailAsk string) ([]*ent.Agent, error) {
	// проверим, есть ли group с таким именем
	g, err := u.repoGroup.GetGroup(ctx, groupName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, me.ErrGroupNotExist
		}
		return nil, err
	}
	if emailAsk == viper.GetString("root_email") {
		return u.repoPrivelege.GetAgents(ctx, g.ID)
	}
	// проверим, существует ли пользователь
	uDB, err := u.repoUser.GetByEmail(ctx, emailAsk)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, me.ErrUserNotExist
		}
		return nil, err
	}
	// удостоверимся, что пользователь владелец группы
	_, err = u.repoGroup.IsOwnerOfGroup(ctx, uDB.ID, groupName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, me.ErrUserIsNotOwner
		}
		return nil, err
	}
	return u.repoPrivelege.GetAgents(ctx, g.ID)
}

func (u *UsecaseLayer) CanExecute(ctx context.Context, userEmail, agentName string) (bool, error) {
	// проверяем, существует ли пользователь, права которого хотим проверить
	uDB, err := u.repoUser.GetByEmail(ctx, userEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, me.ErrUserNotExist
		}
		return false, err
	}
	// проверяем, существует ли агент, к которому хочет обратиться пользователь
	_, err = u.repoAgent.Read(ctx, agentName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, me.ErrAgentNotExist
		}
		return false, err
	}
	// получим список групп, в которые пользователь входит
	// и после для каждой группы получим список доступных агентов
	// если хоть в одном из них окажется agentName, то пользователь имеет доступ к агенту
	groups, err := u.repoGroup.GetUserGroups(ctx, uDB.ID)
	if err != nil {
		// пользователь не находится в какой-либо группе
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	// ключ - название сервиса, значение - пустая структура, которая не занимает память
	userAvaliableAgents := make(map[string]struct{})
	for _, g := range groups {
		as, err := u.repoPrivelege.GetAgents(ctx, g.ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			return false, err
		}
		for _, a := range as {
			userAvaliableAgents[a.Name] = struct{}{}
		}
	}
	if _, ok := userAvaliableAgents[agentName]; !ok {
		return false, nil
	}
	return true, nil
}
