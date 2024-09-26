package user

import (
	"context"
	"fmt"

	ent "github.com/cantylv/authorization-service/internal/entity"
	"github.com/cantylv/authorization-service/internal/utils/myerrors"
	"github.com/jackc/pgx/v5"
)

//go:generate mockgen -source ./repo.go -destination=./mocks/repo.go -package=mock_repo
type Repo interface {
	GetByEmail(ctx context.Context, email string) (*ent.User, error)
	DeleteByEmail(ctx context.Context, email string) error
	Create(ctx context.Context, initData *ent.User) (*ent.User, error)
	// Update(ctx context.Context, updateData *ent.User) (*ent.User, error)
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
	sqlRowGetByEmail = fmt.Sprintf(
		`SELECT %s FROM user WHERE email=$1`,
		user_fields,
	)

	sqlRowCreateUser = fmt.Sprintf(`
		INSERT INTO "user" (
			email,  
			password,
			first_name,
			last_name,    
		) VALUES ($1, $2, $3, $4) RETURNING %s
	`, user_fields)
)

func (r *RepoLayer) GetByEmail(ctx context.Context, email string) (*ent.User, error) {
	row := r.dbConn.QueryRow(ctx, sqlRowGetByEmail, email)
	var u ent.User
	err := row.Scan(&u.Id, &u.Email, &u.Password, &u.FirstName, &u.LastName)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *RepoLayer) DeleteByEmail(ctx context.Context, email string) error {
	row, err := r.dbConn.Exec(ctx, `DELETE FROM "user" WHERE email = $1`, email)
	if err != nil {
		return err
	}
	if row.RowsAffected() == 0 {
		return myerrors.ErrNoRowsAffected
	}
	return nil
}

func (r *RepoLayer) Create(ctx context.Context, initData *ent.User) (*ent.User, error) {
	row := r.dbConn.QueryRow(ctx, sqlRowCreateUser,
		initData.Email,
		initData.Password,
		initData.FirstName,
		initData.LastName,
	)
	var u ent.User
	err := row.Scan(&u.Id, &u.Email, &u.Password, &u.FirstName, &u.LastName)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// func (repo *RepoLayer) Update(ctx context.Context, updateData *ent.User) (*ent.User, error) {
// 	timeNow := time.Now().UTC().Format(cnst.Timestamptz)
// 	userData := data.GetUpdateInfo()
// 	timeNowMetric := time.Now()
// 	row, err := repo.database.ExecContext(ctx,
// 		`UPDATE "user" SET name = $1, email = $2, phone = $3, img_url = $4, password = $5, card_number = $6,
//                   address = $7, updated_at = $8 WHERE email = $9`,
// 		userData.GetName(), userData.GetEmail(), functions.MaybeNullString(userData.GetPhone()), userData.GetImgUrl(),
// 		userData.GetPassword(), functions.MaybeNullString(userData.GetCardNumber()), functions.MaybeNullString(userData.GetAddress()), timeNow, data.GetEmail())

// 	timeEnd := time.Since(timeNowMetric)
// 	repo.metrics.DatabaseDuration.WithLabelValues(metrics.UPDATE).Observe(float64(timeEnd.Milliseconds()))
// 	if err != nil {
// 		return err
// 	}
// 	numRows, err := row.RowsAffected()
// 	if err != nil {
// 		return err
// 	}
// 	if numRows == 0 {
// 		return myerrors.SqlNoRowsUserRelation
// 	}
// 	return nil
// }
