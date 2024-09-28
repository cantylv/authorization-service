package config

import (
	"fmt"
	"os"
	"time"

	"github.com/satori/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// setDefault устанавливает переменные конфигурации viper по умолчанию. Используется для случая, 
// когда файл конфигурации не был найден. 
func setDefault() {
	// PROJECT
	if secretKey := os.Getenv("SECRET_KEY"); secretKey == "" {
		viper.SetDefault("secret_key", uuid.NewV4().String())
	}
	// POSTGRES
	if port := os.Getenv("POSTGRES_PORT"); port == "" {
		viper.SetDefault("postgres.port", 5432)
	}
	if host := os.Getenv("POSTGRES_CONNECTION_HOST"); host == "" {
		viper.SetDefault("postgres.connectionHost", "localhost")
	}
	viper.SetDefault("postgres.sslmode", false)
	// SERVER
	if address := os.Getenv("SERVER_ADDRESS"); address == "" {
		viper.SetDefault("server.address", "localhost:8000")
	}
	if writeTimeout := os.Getenv("SERVER_WRITE_TIMEOUT"); writeTimeout == "" {
		viper.SetDefault("server.write_timeout", 5*time.Second)
	}
	if readTimeout := os.Getenv("SERVER_READ_TIMEOUT"); readTimeout == "" {
		viper.SetDefault("server.read_timeout", 5*time.Second)
	}
	if idleTimeout := os.Getenv("SERVER_IDLE_TIMEOUT"); idleTimeout == "" {
		viper.SetDefault("server.idle_timeout", 3*time.Second)
	}
	if shutdownDuration := os.Getenv("SERVER_SHUTDOWN_DURATION"); shutdownDuration == "" {
		viper.SetDefault("server.shutdown_duration", 10*time.Second)
	}
}

// Read получает переменные из среды и файла конфигурации
func Read(configFilePath string, logger *zap.Logger) {
	setDefault()
	viper.SetConfigFile(configFilePath)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(*os.PathError); !ok {
			logger.Fatal(fmt.Sprintf("fatal error config file: %v", err))
		}
		logger.Warn(fmt.Sprintf("configuration file is not found, programm will be executed within default configuration: %v", err))
		return
	}
	logger.Info("successful read of configuration")
}
