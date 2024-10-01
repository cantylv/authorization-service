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
	Create(ctx context.Context, groupID, agentID int) (*ent.Privelege, error)
	Delete(ctx context.Context, groupID, agentID int) error
	GetAgents(ctx context.Context, groupID int) ([]*ent.Agent, error)
	// CheckPrivelege(ctx context.Context, userID string, privelegeID int) error
	// SetPrivelege(ctx context.Context, userID string, privelegeID int) (*ent.Role, error)
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
	sqlRowCreatePrivelege = `
		INSERT INTO privelege(group_id, agent_id) 
		VALUES ($1, $2) RETURNING id, group_id, agent_id
	`
	sqlRowGetAgents = `
		SELECT a.id, a.name 
		FROM agent a
		JOIN privelege p ON a.id = p.agent_id
		WHERE p.group_id = $1
	`
)

func (r *RepoLayer) Create(ctx context.Context, groupID, agentID int) (*ent.Privelege, error) {
	row := r.dbConn.QueryRow(ctx, sqlRowCreatePrivelege, groupID, agentID)
	var rl ent.Privelege
	err := row.Scan(&rl.ID, &rl.GroupID, &rl.AgentID)
	if err != nil {
		return nil, err
	}
	return &rl, nil
}

func (r *RepoLayer) Delete(ctx context.Context, groupID, agentID int) error {
	tag, err := r.dbConn.Exec(ctx, `DELETE FROM privelege WHERE group_id=$1 AND agent_id=$2`, groupID, agentID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return me.ErrNoRowsAffected
	}
	return nil
}

func (r *RepoLayer) GetAgents(ctx context.Context, groupID int) ([]*ent.Agent, error) {
	rows, err := r.dbConn.Query(ctx, sqlRowGetAgents, groupID)
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

// func (r *RepoLayer) CheckPrivelege(ctx context.Context, userID string, privelegeID int) error {
// 	tag, err := r.dbConn.Exec(ctx, `SELECT 1 FROM role WHERE user_id=$1 AND privelege_id=$2`, userID, privelegeID)
// 	if err != nil {
// 		return err
// 	}
// 	if tag.RowsAffected() == 0 {
// 		return me.ErrNoRowsAffected
// 	}
// 	return nil
// }

// func (r *RepoLayer) SetPrivelege(ctx context.Context, userID string, privelegeID int) (*ent.Role, error) {
// 	row := r.dbConn.QueryRow(ctx, sqlRowSetPrivelege, userID, privelegeID)
// 	var rl ent.Role
// 	err := row.Scan(&rl.ID, &rl.UserID, &rl.PrivelegeID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &rl, nil
// }
