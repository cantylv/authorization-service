package privelege

import (
	"context"
	"database/sql"
	"errors"

	ent "github.com/cantylv/authorization-service/internal/entity"
	me "github.com/cantylv/authorization-service/internal/utils/myerrors"
	"github.com/jackc/pgx/v5"
)

type Repo interface {
	CreateGroupAgent(ctx context.Context, groupID, agentID int) (*ent.GroupPrivelege, error)
	CreateUserAgent(ctx context.Context, userID string, agentID int) (*ent.UserPrivelege, error)
	DeleteGroupAgent(ctx context.Context, groupID, agentID int) error
	DeleteUserAgent(ctx context.Context, userID string, agentID int) error
	GetGroupAgents(ctx context.Context, groupID int) ([]*ent.Agent, error)
	GetUserAgents(ctx context.Context, userID string) ([]*ent.Agent, error)
}

var _ Repo = (*RepoLayer)(nil)

type RepoLayer struct {
	dbConn *pgx.Conn
}

func NewRepoLayer(dbConn *pgx.Conn) *RepoLayer {
	return &RepoLayer{
		dbConn: dbConn,
	}
}

var (
	sqlRowCreateGroupPrivelege = `
		INSERT INTO privelege_group(group_id, agent_id) 
		VALUES ($1, $2) RETURNING id, group_id, agent_id
	`
	sqlRowCreateUserPrivelege = `
		INSERT INTO privelege_user(user_id, agent_id) 
		VALUES ($1, $2) RETURNING id, user_id, agent_id
	`
	sqlRowGetGroupAgents = `
		SELECT a.id, a.name 
		FROM agent a
		JOIN privelege_group p ON a.id = p.agent_id
		WHERE p.group_id = $1
	`
	sqlRowGetUserAgents = `
		SELECT a.id, a.name 
		FROM agent a
		JOIN privelege_user p ON a.id = p.agent_id
		WHERE p.user_id = $1
	`
)

func (r *RepoLayer) CreateGroupAgent(ctx context.Context, groupID, agentID int) (*ent.GroupPrivelege, error) {
	row := r.dbConn.QueryRow(ctx, sqlRowCreateGroupPrivelege, groupID, agentID)
	var rl ent.GroupPrivelege
	err := row.Scan(&rl.ID, &rl.GroupID, &rl.AgentID)
	if err != nil {
		return nil, err
	}
	return &rl, nil
}

func (r *RepoLayer) CreateUserAgent(ctx context.Context, userID string, agentID int) (*ent.UserPrivelege, error) {
	row := r.dbConn.QueryRow(ctx, sqlRowCreateUserPrivelege, userID, agentID)
	var rl ent.UserPrivelege
	err := row.Scan(&rl.ID, &rl.UserID, &rl.AgentID)
	if err != nil {
		return nil, err
	}
	return &rl, nil
}

func (r *RepoLayer) DeleteGroupAgent(ctx context.Context, groupID, agentID int) error {
	tag, err := r.dbConn.Exec(ctx, `DELETE FROM privelege_group WHERE group_id=$1 AND agent_id=$2`, groupID, agentID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return me.ErrNoRowsAffected
	}
	return nil
}

func (r *RepoLayer) DeleteUserAgent(ctx context.Context, userID string, agentID int) error {
	tag, err := r.dbConn.Exec(ctx, `DELETE FROM privelege_user WHERE user_id=$1 AND agent_id=$2`, userID, agentID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return me.ErrNoRowsAffected
	}
	return nil
}

func (r *RepoLayer) GetGroupAgents(ctx context.Context, groupID int) ([]*ent.Agent, error) {
	rows, err := r.dbConn.Query(ctx, sqlRowGetGroupAgents, groupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	var agents []*ent.Agent
	for rows.Next() {
		var a ent.Agent
		err = rows.Scan(&a.ID, &a.Name)
		if err != nil {
			return nil, err
		}
		agents = append(agents, &a)
	}
	return agents, nil
}

func (r *RepoLayer) GetUserAgents(ctx context.Context, userID string) ([]*ent.Agent, error) {
	rows, err := r.dbConn.Query(ctx, sqlRowGetUserAgents, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	var agents []*ent.Agent
	for rows.Next() {
		var a ent.Agent
		err = rows.Scan(&a.ID, &a.Name)
		if err != nil {
			return nil, err
		}
		agents = append(agents, &a)
	}
	return agents, nil
}
