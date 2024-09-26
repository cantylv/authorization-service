package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

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
	logger.Info("postgresql connected successfully")
	return conn
}
