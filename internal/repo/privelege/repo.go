package privelege

import (
	"context"

	ent "github.com/cantylv/authorization-service/internal/entity"
	"github.com/jackc/pgx/v5"
)

type Repo interface {
	// Create(ctx context.Context, ownerID string, privelegeID int) (*ent.Role, error)
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

// var (
// 	sqlRowCreatePrivelege = `INSERT INTO role(user_id, privelege_id) VALUES ($1, $2) RETURNING id, user_id, privelege_id`
// 	sqlRowSetPrivelege    = `INSERT INTO role(user_id, privelege_id) VALUES($1, $2) RETURNING id, user_id, privelege_id`
// )

func (r *RepoLayer) CreateAgent(ctx context.Context, name string) (*ent.Agent, error) {
	row := r.dbConn.QueryRow(ctx, `INSERT INTO agent(name) VALUES($1) RETURNING id, name`, name)
	var a ent.Agent
	err := row.Scan(&a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// func (r *RepoLayer) Create(ctx context.Context, ownerID string, privelegeID int) (*ent.Role, error) {
// 	row := r.dbConn.QueryRow(ctx, sqlRowCreatePrivelege, ownerID, privelegeID)
// 	var rl ent.Role
// 	err := row.Scan(&rl.ID, &rl.UserID, &rl.PrivelegeID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &rl, nil
// }

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
