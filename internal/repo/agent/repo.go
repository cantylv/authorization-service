package agent

import (
	"context"

	ent "github.com/cantylv/authorization-service/internal/entity"
	me "github.com/cantylv/authorization-service/internal/utils/myerrors"
	"github.com/jackc/pgx/v5"
)

type Repo interface {
	Read(ctx context.Context, name string) (*ent.Agent, error)
	GetAll(ctx context.Context) ([]*ent.Agent, error)
	Create(ctx context.Context, name string) (*ent.Agent, error)
	Delete(ctx context.Context, id int) error
	IsGroupAgent(ctx context.Context, groupID, agentID int) (bool, error)
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

func (r *RepoLayer) Read(ctx context.Context, name string) (*ent.Agent, error) {
	row := r.dbConn.QueryRow(ctx, `SELECT id, name FROM agent WHERE name=$1`, name)
	var a ent.Agent
	err := row.Scan(&a.ID, &a.Name)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *RepoLayer) GetAll(ctx context.Context) ([]*ent.Agent, error) {
	rows, err := r.dbConn.Query(ctx, `SELECT id, name FROM agent`)
	if err != nil {
		return nil, err
	}
	var as []*ent.Agent
	for rows.Next() {
		var a ent.Agent
		err := rows.Scan(&a.ID, &a.Name)
		if err != nil {
			return nil, err
		}
		as = append(as, &a)
	}
	return as, nil
}

func (r *RepoLayer) Create(ctx context.Context, name string) (*ent.Agent, error) {
	row := r.dbConn.QueryRow(ctx, `INSERT INTO agent(name) VALUES($1) RETURNING id, name`, name)
	var a ent.Agent
	err := row.Scan(&a.ID, &a.Name)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *RepoLayer) Delete(ctx context.Context, id int) error {
	tag, err := r.dbConn.Exec(ctx, `DELETE FROM agent WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return me.ErrNoRowsAffected
	}
	return nil
}

func (r *RepoLayer) IsGroupAgent(ctx context.Context, groupID, agentID int) (bool, error) {
	row := r.dbConn.QueryRow(ctx, `SELECT 1 FROM privelege WHERE group_id=$1 AND agent_id=$2`, groupID, agentID)
	var isGroupAgent int
	err := row.Scan(&isGroupAgent)
	if err != nil {
		return false, err
	}
	return true, nil
}
