package group

import (
	"context"

	ent "github.com/cantylv/authorization-service/internal/entity"
	"github.com/jackc/pgx/v5"
)

type Repo interface {
	Create(ctx context.Context, ownerID string) (*ent.Group, error)
	UpdateOwner(ctx context.Context, groupID int, newOwnerID string) (*ent.Group, error)
	GetParticipants(ctx context.Context, groupID int) ([]*ent.User, error)
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
	user_fields = "id, email, password, first_name, last_name"
)

var (
	sqlRowGetParticipants = `
		SELECT u.id, u.email, u.first_name, u.last_name
		FROM "user" u
		JOIN participant p ON u.id = p.user_id
		WHERE p.group_id = $1
	`
)

func (r *RepoLayer) Create(ctx context.Context, ownerID string) (*ent.Group, error) {
	row := r.dbConn.QueryRow(ctx, `INSERT INTO group(owner_id) VALUES ($1) RETURNING id, owner_id`, ownerID)
	var g ent.Group
	err := row.Scan(&g.ID, &g.OwnerID)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *RepoLayer) UpdateOwner(ctx context.Context, groupID int, newOwnerID string) (*ent.Group, error) {
	row := r.dbConn.QueryRow(ctx, `UPDATE group SET owner_id=$1 WHERE id=$2`, newOwnerID, groupID)
	var g ent.Group
	err := row.Scan(&g.ID, &g.OwnerID)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *RepoLayer) GetParticipants(ctx context.Context, groupID int) ([]*ent.User, error) {
	rows, err := r.dbConn.Query(ctx, sqlRowGetParticipants, groupID)
	if err != nil {
		return nil, err
	}
	var us []*ent.User
	for rows.Next() {
		var u ent.User
		err := rows.Scan(&u.Id, &u.Email, &u.FirstName, &u.LastName)
		if err != nil {
			return nil, err
		}
		us = append(us, &u)
	}
	return us, nil
}
