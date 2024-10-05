package archive

import (
	"context"

	ent "github.com/cantylv/authorization-service/microservices/archive_manager/internal/entity"
	"github.com/jackc/pgx/v5"
)

type Repo interface {
	Get(ctx context.Context) ([]*ent.Record, error)
}

var _ Repo = (*RepoLayer)(nil)

type RepoLayer struct {
	dbconn *pgx.Conn
}

func NewRepoLayer(conn *pgx.Conn) *RepoLayer {
	return &RepoLayer{
		dbconn: conn,
	}
}

func (r *RepoLayer) Get(ctx context.Context) ([]*ent.Record, error) {
	rows, err := r.dbconn.Query(ctx, `SELECT id, text FROM record`)
	if err != nil {
		return nil, err
	}
	var records []*ent.Record
	for rows.Next() {
		var rec ent.Record
		err := rows.Scan(&rec.ID, &rec.Text)
		if err != nil {
			return nil, err
		}
		records = append(records, &rec)
	}
	return records, nil
}
