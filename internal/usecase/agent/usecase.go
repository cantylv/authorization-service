package agent

import (
	"context"
	"database/sql"
	"errors"

	ent "github.com/cantylv/authorization-service/internal/entity"
	"github.com/cantylv/authorization-service/internal/repo/agent"
	me "github.com/cantylv/authorization-service/internal/utils/myerrors"
	"github.com/spf13/viper"
)

type Usecase interface {
	CreateAgent(ctx context.Context, emailCreator, agentName string) (*ent.Agent, error)
	DeleteAgent(ctx context.Context, emailCreator, agentName string) error
	GetAgents(ctx context.Context, emailCreator string) ([]*ent.Agent, error)
}

var _ Usecase = (*UsecaseLayer)(nil)

type UsecaseLayer struct {
	repoAgent agent.Repo
}

// NewUsecaseLayer возвращает структуру уровня usecase, управляющую агентами серверной архитектуры
func NewUsecaseLayer(repoAgent agent.Repo) *UsecaseLayer {
	return &UsecaseLayer{
		repoAgent: repoAgent,
	}
}

// CreateAgent создает агента, его создать может только root пользователь
func (u *UsecaseLayer) CreateAgent(ctx context.Context, emailCreator, agentName string) (*ent.Agent, error) {
	if emailCreator != viper.GetString("root_email") {
		return nil, me.ErrOnlyRootCanAddAgent
	}
	// проверяем, есть ли уже агент с таким именем, если есть, то возвращаем ошибку
	a, err := u.repoAgent.Read(ctx, agentName)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if a != nil {
		return nil, me.ErrAgentAlreadyExist
	}
	// создаем
	a, err = u.repoAgent.Create(ctx, agentName)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// DeleteAgent удаляет агента, его удалить может только root пользователь
func (u *UsecaseLayer) DeleteAgent(ctx context.Context, emailCreator, agentName string) error {
	if emailCreator != viper.GetString("root_email") {
		return me.ErrOnlyRootCanDeleteAgent
	}
	// проверяем, есть ли агент с таким именем
	a, err := u.repoAgent.Read(ctx, agentName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return me.ErrAgentNotExist
		}
		return err
	}
	return u.repoAgent.Delete(ctx, a.ID)
}

// DeleteAgent удаляет агента, его удалить может только root пользователь
func (u *UsecaseLayer) GetAgents(ctx context.Context, emailCreator string) ([]*ent.Agent, error) {
	if emailCreator != viper.GetString("root_email") {
		return nil, me.ErrOnlyRootCanGetAgents
	}
	// проверяем, есть ли агент с таким именем
	a, err := u.repoAgent.GetAll(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return a, nil
}
