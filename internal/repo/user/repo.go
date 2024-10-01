package user

import (
	"context"
	"fmt"

	ent "github.com/cantylv/authorization-service/internal/entity"
	me "github.com/cantylv/authorization-service/internal/utils/myerrors"
	"github.com/jackc/pgx/v5"
)

//go:generate mockgen -source ./repo.go -destination=./mocks/repo.go -package=mock_repo
type Repo interface {
	GetByEmail(ctx context.Context, email string) (*ent.User, error)
	DeleteByEmail(ctx context.Context, email string) error
	Create(ctx context.Context, initData *ent.User) (*ent.User, error)
}

var _ Repo = (*RepoLayer)(nil)

type RepoLayer struct {
	dbConn *pgx.Conn
}

// NewRepoLayer возвращает структуру уровня repository. Позволяет работать с пользователем (crd).
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
		`SELECT %s FROM "user" WHERE email=$1`,
		user_fields,
	)
	sqlRowCreateUser = fmt.Sprintf(`
		INSERT INTO "user" (
			email,  
			password,
			first_name,
			last_name    
		) VALUES ($1, $2, $3, $4) RETURNING %s`, user_fields)
)

// GetByEmail позволяет получить пользователя
func (r *RepoLayer) GetByEmail(ctx context.Context, email string) (*ent.User, error) {
	row := r.dbConn.QueryRow(ctx, sqlRowGetByEmail, email)
	var u ent.User
	err := row.Scan(&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// DeleteByEmail позволяет удалить пользователя из системы. Пользователя можно удалить только
// в том случае, если он не является ответственным за какую-либо группу. Если он таковым является,
// необходимо сперва поменять ответственного. Это сделать может только root.
func (r *RepoLayer) DeleteByEmail(ctx context.Context, email string) error {
	row, err := r.dbConn.Exec(ctx, `DELETE FROM "user" WHERE email = $1`, email)
	if err != nil {
		return err
	}
	if row.RowsAffected() == 0 {
		return me.ErrNoRowsAffected
	}
	return nil
}

// Create позволяет создать пользователя. В процессе выполнения добавляет пользователя в группу 'users'.
// Каждый пользователь в системе принадлежит базовой группе 'users'.
func (r *RepoLayer) Create(ctx context.Context, initData *ent.User) (*ent.User, error) {
	tx, err := r.dbConn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// создаем пользователя
	rowUser := tx.QueryRow(ctx, sqlRowCreateUser,
		initData.Email,
		initData.Password,
		initData.FirstName,
		initData.LastName,
	)
	var u ent.User
	err = rowUser.Scan(&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName)
	if err != nil {
		return nil, err
	}

	// получим id группы, чтобы потом записать пользователя в группу
	var groupID int
	rowGroup := tx.QueryRow(ctx, `SELECT id FROM "group" WHERE name='users'`)
	err = rowGroup.Scan(&groupID)
	if err != nil {
		return nil, err
	}
	// теперь нужно записать пользователя в группу пользователей
	tag, err := tx.Exec(ctx, `INSERT INTO participation(user_id, group_id) VALUES ($1, $2)`, u.ID, groupID)
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
	return &u, nil
}
