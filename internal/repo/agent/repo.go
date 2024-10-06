package agent

import (
	"context"

	ent "github.com/cantylv/authorization-service/internal/entity"
	me "github.com/cantylv/authorization-service/internal/utils/myerrors"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
)

type Repo interface {
	Read(ctx context.Context, name string) (*ent.Agent, error)
	GetAll(ctx context.Context) ([]*ent.Agent, error)
	Create(ctx context.Context, name string) (*ent.Agent, error)
	Delete(ctx context.Context, id int) error
	IsGroupAgent(ctx context.Context, groupID, agentID int) (bool, error)
	IsUserAgent(ctx context.Context, userID string, agentID int) (bool, error)
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
	tx, err := r.dbConn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()
	row := tx.QueryRow(ctx, `INSERT INTO agent(name) VALUES($1) RETURNING id, name`, name)
	var a ent.Agent
	err = row.Scan(&a.ID, &a.Name)
	if err != nil {
		return nil, err
	}
	// нужно добавить root в права на владение агентом
	var userID string
	row = tx.QueryRow(ctx, `SELECT id FROM "user" where email=$1`, viper.GetString("root_email"))
	err = row.Scan(&userID)
	if err != nil {
		return nil, err
	}
	tag, err := tx.Exec(ctx, `INSERT INTO privelege_user(agent_id, user_id) VALUES($1, $2)`, a.ID, userID)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, me.ErrNoRowsAffected
	}
	// если все прошло успешно, коммитим транзакцию
	if err = tx.Commit(ctx); err != nil {
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
	row := r.dbConn.QueryRow(ctx, `SELECT 1 FROM privelege_group WHERE group_id=$1 AND agent_id=$2`, groupID, agentID)
	var isGroupAgent int
	err := row.Scan(&isGroupAgent)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *RepoLayer) IsUserAgent(ctx context.Context, userID string, agentID int) (bool, error) {
	row := r.dbConn.QueryRow(ctx, `SELECT 1 FROM privelege_user WHERE user_id=$1 AND agent_id=$2`, userID, agentID)
	var isUserAgent int
	err := row.Scan(&isUserAgent)
	if err != nil {
		return false, err
	}
	return true, nil
}
