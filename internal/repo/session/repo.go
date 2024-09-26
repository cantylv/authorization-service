package session

import (
	"context"
	"fmt"
	"time"

	ent "github.com/cantylv/authorization-service/internal/entity"
	mc "github.com/cantylv/authorization-service/internal/utils/myconstants"
	me "github.com/cantylv/authorization-service/internal/utils/myerrors"
	"github.com/jackc/pgx/v5"
	"github.com/satori/uuid"
)

type Repo interface {
	CreateSession(ctx context.Context, userID string, initData *ent.Session) (*ent.Session, error)
	ReadSessions(ctx context.Context, userID string) ([]*ent.Session, error)
	DeleteUserSessions(ctx context.Context, userID string) error
	DeleteSessionByToken(ctx context.Context, refreshToken string) error
}

var _ Repo = (*RepoLayer)(nil)

type RepoLayer struct {
	dbConn *pgx.Conn
}

// NewRepoLayer возвращает указатель на структуру репозитория с сессиями (доступно: создание, получение, удаление).
func NewRepoLayer(dbConn *pgx.Conn) *RepoLayer {
	return &RepoLayer{
		dbConn: dbConn,
	}
}

var (
	tokenData = "id, user_id, refresh_token, fingerprint, user_ip_address, expires_at"
)

var (
	sqlRowCreateSession = fmt.Sprintf(`
	INSERT INTO session (
		user_id,
		refresh_token,
		fingerprint,
		user_ip_address,
		expires_at
	) VALUES ($1, $2, $3, $4, $5) RETURNING %s`, tokenData)
	sqlRowReadSessions = fmt.Sprintf(`SELECT %s FROM TABLE session WHERE user_id=$1`, tokenData)
)

func (r *RepoLayer) CreateSession(ctx context.Context, userID string, initData *ent.Session) (*ent.Session, error) {
	row := r.dbConn.QueryRow(ctx, sqlRowCreateSession,
		userID,
		uuid.NewV4().String(), // refresh_token
		initData.Fingerprint,
		initData.UserIpAddress,
		time.Now().AddDate(0, 0, mc.DayExpRefreshToken), // expires_at
	)
	var t ent.Session
	err := row.Scan(&t.ID, &t.UserID, &t.RefreshToken, &t.Fingerprint, &t.UserIpAddress, &t.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *RepoLayer) ReadSessions(ctx context.Context, userID string) ([]*ent.Session, error) {
	rows, err := r.dbConn.Query(ctx, sqlRowReadSessions, userID)
	if err != nil {
		return nil, err
	}
	var ts []*ent.Session
	for rows.Next() {
		var t ent.Session
		err := rows.Scan(&t.ID, &t.UserID, &t.RefreshToken, &t.Fingerprint, &t.UserIpAddress, &t.ExpiresAt)
		if err != nil {
			return nil, err
		}
		ts = append(ts, &t)
	}
	return ts, nil
}

func (r *RepoLayer) DeleteUserSessions(ctx context.Context, userID string) error {
	tag, err := r.dbConn.Exec(ctx, `DELETE FROM TABLE session WHERE user_id=$1`, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return me.ErrNoRowsAffected
	}
	return nil
}

func (r *RepoLayer) DeleteSessionByToken(ctx context.Context, refreshToken string) error {
	tag, err := r.dbConn.Exec(ctx, `DELETE FROM TABLE session WHERE refresh_token=$1`, refreshToken)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return me.ErrNoRowsAffected
	}
	return nil
}
