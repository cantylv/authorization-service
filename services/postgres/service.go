package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/cantylv/authorization-service/internal/entity/dto"
	"github.com/cantylv/authorization-service/internal/utils/functions"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Init инициализирует клиента PostgreSQL.
func Init(logger *zap.Logger) *pgx.Conn {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		viper.GetString("postgres.user"),
		viper.GetString("postgres.password"),
		viper.GetString("postgres.connectionHost"),
		viper.GetInt("postgres.port"),
		viper.GetString("postgres.database_name"),
		viper.GetString("postgres.sslmode"),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		logger.Fatal(fmt.Sprintf("error while connecting to postgresql: %v", err))
	}
	const maxPingAttempts = 3

	var successConn bool
	for i := 0; i < maxPingAttempts; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		err := conn.Ping(ctx)
		cancel()
		if err == nil {
			successConn = true
			break
		}
		logger.Warn(fmt.Sprintf("error while ping to postgresql: %v", err))
	}
	if !successConn {
		logger.Fatal("can't establish connection to postgresql")
	}

	err = createRootUser(conn)
	if err != nil {
		logger.Fatal(fmt.Sprintf("error while creating root user: %v", err))
	}

	logger.Info("postgresql connected successfully")
	return conn
}

func isExistRootUser(conn *pgx.Conn) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row := conn.QueryRow(ctx, `SELECT 1 FROM "user" WHERE email=$1`, viper.GetString("root_email"))
	var exist int
	err := row.Scan(&exist)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func createRootUser(conn *pgx.Conn) error {
	isExist, err := isExistRootUser(conn)
	if err != nil {
		return err
	}
	if isExist {
		return nil
	}
	// создаем root пользователя
	rootUser := dto.CreateData{
		Email:     viper.GetString("root_email"),
		Password:  viper.GetString("root_password"),
		FirstName: viper.GetString("root_first_name"),
		LastName:  viper.GetString("root_last_name"),
	}
	err = rootUser.Validate()
	if err != nil {
		return err
	}
	hashedRootPassword, err := functions.GetHashedPassword(viper.GetString("root_password"))
	if err != nil {
		return err
	}
	rootUser.Password = hashedRootPassword
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row := conn.QueryRow(ctx,
		`INSERT INTO "user"(email, password, first_name, last_name) VALUES ($1, $2, $3, $4) RETURNING id`,
		rootUser.Email, rootUser.Password, rootUser.FirstName, rootUser.LastName)
	var userID string
	err = row.Scan(&userID)
	if err != nil {
		return err
	}
	// создаем группу обычных пользователей
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row = conn.QueryRow(ctx, `INSERT INTO "group"(name, owner_id) VALUES('users', $1) RETURNING id`, userID)
	var groupID int
	err = row.Scan(&groupID)
	if err != nil {
		return err
	}
	// добавляем root пользователя в группу обычных пользователей
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tag, err := conn.Exec(ctx, `INSERT INTO participation(user_id, group_id) VALUES($1, $2)`, userID, groupID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("root user was not added to group 'users'")
	}
	return nil
}
