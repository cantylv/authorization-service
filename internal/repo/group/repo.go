package group

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	ent "github.com/cantylv/authorization-service/internal/entity"
	"github.com/cantylv/authorization-service/internal/entity/dto"
	me "github.com/cantylv/authorization-service/internal/utils/myerrors"
	"github.com/jackc/pgx/v5"
)

type Repo interface {
	GetGroup(ctx context.Context, groupName string) (*ent.Group, error)
	GetBid(ctx context.Context, userID, groupName string) (*dto.Bid, error)
	AddUserToGroup(ctx context.Context, userID string, groupID int) error
	ApproveGroupCreation(ctx context.Context, ownerID, rootUserID, groupName string) (*ent.Group, error)
	RejectGroupCreation(ctx context.Context, groupName, userID string) (*dto.Bid, error)
	MakeBidGroupCreation(ctx context.Context, ownerID, groupName string) (*dto.Bid, error)
	IsParticipantOfGroup(ctx context.Context, userID string, groupID int) (bool, error)
	IsOwnerOfGroup(ctx context.Context, userID, groupName string) (bool, error)
	GetCommonGroups(ctx context.Context, userID1, userID2 string) ([]*ent.Group, error)
	GetUserGroups(ctx context.Context, userID string) ([]*ent.Group, error)
	KickUserFromGroup(ctx context.Context, userID string, groupID int) error
	CreateGroup(ctx context.Context, userID, groupName string) (*ent.Group, error)
	OwnerGroups(ctx context.Context, userID string) ([]*ent.Group, error)
	UpdateOwner(ctx context.Context, groupID int, newOwnerID string) (*ent.Group, error)
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
	group_fiels = "id, name, owner_id"
)

var (
	sqlRowGetParticipants = `
		SELECT u.id, u.email, u.first_name, u.last_name
		FROM "user" u
		JOIN participant p ON u.id = p.user_id
		WHERE p.group_id = $1
	`
	sqlRowCreateGroup = fmt.Sprintf(`INSERT INTO "group"(name, owner_id) VALUES ($1, $2) RETURNING %s`, group_fiels)
	sqlRowMakeBid     = `
		INSERT INTO bid (group_name, user_id, status) 
		VALUES ($1, $2, $3) 
		RETURNING id, group_name, user_id, status
	`
	sqlRowRejectBid = `
		UPDATE bid 
		SET status = 'rejected' 
		WHERE group_name = $1 AND user_id = $2
		RETURNING id, group_name, user_id, status
	`
	sqlRowGetOwnerGroups = `
		SELECT g.id, g.name, g.owner_id
		FROM "group" g
		JOIN "user" u ON g.owner_id = u.id
		WHERE u.id = $1
	`
	sqlRowGetBid = `
		SELECT id, group_name, user_id, status 
		FROM bid 
		WHERE user_id=$1 AND group_name=$2 AND status != 'rejected'
	`
	sqlRowAddUserToGroup = `INSERT INTO participation(user_id, group_id) VALUES ($1, $2)`
)

// GetGroup возвращает данные о группе
func (r *RepoLayer) GetGroup(ctx context.Context, groupName string) (*ent.Group, error) {
	row := r.dbConn.QueryRow(ctx, `SELECT id, name, owner_id FROM "group" WHERE name=$1`, groupName)
	var g ent.Group
	err := row.Scan(&g.ID, &g.Name, &g.OwnerID)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

// GetBid возвращает данные об активной заявке
func (r *RepoLayer) GetBid(ctx context.Context, userID, groupName string) (*dto.Bid, error) {
	row := r.dbConn.QueryRow(ctx, sqlRowGetBid, userID, groupName)
	var b dto.Bid
	err := row.Scan(&b.ID, &b.GroupName, &b.UserId, &b.Status)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// GetCommonGroups добавляет пользователя в группу
func (r *RepoLayer) AddUserToGroup(ctx context.Context, ownerID string, groupID int) error {
	tag, err := r.dbConn.Exec(ctx, sqlRowAddUserToGroup, ownerID, groupID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return me.ErrNoRowsAffected
	}
	return nil
}

// IsParticipantOfGroup определяет, является ли пользователь членом группы
func (r *RepoLayer) IsParticipantOfGroup(ctx context.Context, userID string, groupID int) (bool, error) {
	row := r.dbConn.QueryRow(ctx, `SELECT 1 FROM participation WHERE user_id=$1 AND group_id=$2`, userID, groupID)
	var isParticipant int
	err := row.Scan(&isParticipant)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *RepoLayer) IsOwnerOfGroup(ctx context.Context, userID, groupName string) (bool, error) {
	row := r.dbConn.QueryRow(ctx, `SELECT 1 FROM "group" WHERE owner_id=$1 AND name=$2`, userID, groupName)
	var isOwner int
	err := row.Scan(&isOwner)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetCommonGroups возвращает список совместных групп двух пользователей
func (r *RepoLayer) GetCommonGroups(ctx context.Context, userID1, userID2 string) ([]*ent.Group, error) {
	rows, err := r.dbConn.Query(ctx, `
		SELECT g.id, g.name, g.owner_id
		FROM "group" g
		JOIN participation p1 ON g.id = p1.group_id
		JOIN participation p2 ON g.id = p2.group_id
		WHERE p1.user_id = $1 AND p2.user_id = $2;
	`, userID1, userID2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*ent.Group
	for rows.Next() {
		var group ent.Group
		if err := rows.Scan(&group.ID, &group.Name, &group.OwnerID); err != nil {
			return nil, err
		}
		groups = append(groups, &group)
	}

	return groups, nil
}

// GetUserGroups возвращает список групп, в которых пользователь состоит
func (r *RepoLayer) GetUserGroups(ctx context.Context, userID string) ([]*ent.Group, error) {
	rows, err := r.dbConn.Query(ctx, `
		SELECT g.id, g.name, g.owner_id
		FROM "group" g
		JOIN participation p ON p.group_id = g.id
		WHERE p.user_id = $1;
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*ent.Group
	for rows.Next() {
		var group ent.Group
		if err := rows.Scan(&group.ID, &group.Name, &group.OwnerID); err != nil {
			return nil, err
		}
		groups = append(groups, &group)
	}

	return groups, nil
}

// KickUserFromGroup удаляет пользователя из группы
func (r *RepoLayer) KickUserFromGroup(ctx context.Context, userID string, groupID int) error {
	tag, err := r.dbConn.Exec(ctx, `DELETE FROM participation WHERE user_id=$1 AND group_id=$2`, userID, groupID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return me.ErrNoRowsAffected
	}
	return nil
}

// CreateGroup создает группу, а также добавляет создателя в группу (т. participation)
func (r *RepoLayer) CreateGroup(ctx context.Context, userID, groupName string) (*ent.Group, error) {
	tx, err := r.dbConn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()
	// создаем группу
	row := tx.QueryRow(ctx, sqlRowCreateGroup, groupName, userID)
	var g ent.Group
	err = row.Scan(&g.ID, &g.Name, &g.OwnerID)
	if err != nil {
		return nil, err
	}
	// добавляем создателя в его группу
	tag, err := tx.Exec(ctx, sqlRowAddUserToGroup, userID, g.ID)
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
	return &g, nil
}

// ApproveGroupCreation метод, который аппрувит статус создания группы каким-либо пользователем.
// Сперва удаляет запись из таблицы заявок, после добавляет группу в таблицу
func (r *RepoLayer) ApproveGroupCreation(ctx context.Context, ownerID, rootUserID, groupName string) (*ent.Group, error) {
	tx, err := r.dbConn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// удаляем запись из таблицы заявок
	tag, err := tx.Exec(ctx, `DELETE FROM bid WHERE group_name=$1 AND user_id=$2`, groupName, ownerID)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		err = me.ErrNoRowsAffected
		return nil, err
	}
	// создаем новую группу
	row := tx.QueryRow(ctx, `INSERT INTO "group"(name, owner_id) VALUES($1, $2) RETURNING id, name, owner_id`, groupName, ownerID)
	var g ent.Group
	err = row.Scan(&g.ID, &g.Name, &g.OwnerID)
	if err != nil {
		return nil, err
	}
	// добавляем создателя в его группу
	tag, err = tx.Exec(ctx, sqlRowAddUserToGroup, ownerID, g.ID)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, me.ErrNoRowsAffected
	}
	// добавляем root пользователя в группу
	tag, err = tx.Exec(ctx, sqlRowAddUserToGroup, rootUserID, g.ID)
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
	return &g, nil
}

// MakeBidGroupCreation создает заявку на создание группы
func (r *RepoLayer) MakeBidGroupCreation(ctx context.Context, ownerID, groupName string) (*dto.Bid, error) {
	row := r.dbConn.QueryRow(ctx, sqlRowMakeBid, groupName, ownerID, "in_progress")
	var g dto.Bid
	err := row.Scan(&g.ID, &g.GroupName, &g.UserId, &g.Status)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *RepoLayer) UpdateOwner(ctx context.Context, groupID int, newOwnerID string) (*ent.Group, error) {
	tx, err := r.dbConn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()
	// проверяем, есть ли пользователь в группе уже
	row := tx.QueryRow(ctx, `SELECT 1 FROM participation WHERE user_id=$1 AND group_id=$2`, newOwnerID, groupID)
	var isParticipant int
	err = row.Scan(&isParticipant)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	// если пользователь не в группе, то надо его добавить
	if isParticipant != 1 {
		tag, err := tx.Exec(ctx, `INSERT INTO participation(user_id, group_id) VALUES($1, $2)`, newOwnerID, groupID)
		if err != nil {
			return nil, err
		}
		if tag.RowsAffected() == 0 {
			return nil, me.ErrNoRowsAffected
		}
	}
	// делаем пользователя ответственным
	row = tx.QueryRow(ctx, `UPDATE "group" SET owner_id=$1 WHERE id=$2 RETURNING id, name, owner_id`, newOwnerID, groupID)
	var g ent.Group
	err = row.Scan(&g.ID, &g.Name, &g.OwnerID)
	if err != nil {
		return nil, err
	}
	// если все прошло успешно, коммитим транзакцию
	if err = tx.Commit(ctx); err != nil {
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
		err := rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName)
		if err != nil {
			return nil, err
		}
		us = append(us, &u)
	}
	return us, nil
}

func (r *RepoLayer) RejectGroupCreation(ctx context.Context, groupName, userID string) (*dto.Bid, error) {
	row := r.dbConn.QueryRow(ctx, sqlRowRejectBid, groupName, userID)
	var g dto.Bid
	err := row.Scan(&g.ID, &g.GroupName, &g.UserId, &g.Status)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *RepoLayer) OwnerGroups(ctx context.Context, userID string) ([]*ent.Group, error) {
	rows, err := r.dbConn.Query(ctx, sqlRowGetOwnerGroups, userID)
	if err != nil {
		return nil, err
	}

	var groups []*ent.Group
	for rows.Next() {
		var g ent.Group
		err := rows.Scan(&g.ID, &g.Name, &g.OwnerID)
		if err != nil {
			return nil, err
		}
		groups = append(groups, &g)
	}
	return groups, nil
}
