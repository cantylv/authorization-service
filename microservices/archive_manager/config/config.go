package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// readEnvAndSetDefault устанавливает переменные конфигурации viper по умолчанию. Используется для случая,
// когда файл конфигурации не был найден. Использует переменные окружения для настройки.
func readEnvAndSetDefault(logger *zap.Logger) {
	// POSTGRES
	if port := os.Getenv("A_POSTGRES_PORT"); port != "" {
		psqlPort, err := strconv.Atoi(port)
		if err != nil {
			logger.Info("you've passed incorrect value of env variable 'A_POSTGRES_PORT', so it will be with default value 5433")
			viper.SetDefault("postgres.port", 5433)
		} else {
			viper.SetDefault("postgres.port", psqlPort)
		}
	} else {
		viper.SetDefault("postgres.port", 5433)
	}

	if host := os.Getenv("A_POSTGRES_CONNECTION_HOST"); host != "" {
		viper.SetDefault("postgres.connectionHost", host)
	} else {
		viper.SetDefault("postgres.connectionHost", "localhost")
	}

	viper.SetDefault("postgres.sslmode", "disable")
	// SERVER
	if address := os.Getenv("AM_SERVER_ADDRESS"); address != "" {
		viper.SetDefault("archive_manager.address", address)
	} else {
		viper.SetDefault("archive_manager.address", "localhost:8011")
	}

	if writeTimeout := os.Getenv("AM_SERVER_WRITE_TIMEOUT"); writeTimeout != "" {
		timeout, err := time.ParseDuration(writeTimeout)
		if err != nil {
			logger.Info("you've passed incorrect value of env variable 'AM_SERVER_WRITE_TIMEOUT', so it will be with default value 5s")
			viper.SetDefault("archive_manager.write_timeout", 5*time.Second)
		} else {
			viper.SetDefault("archive_manager.write_timeout", timeout)
		}
	} else {
		viper.SetDefault("archive_manager.write_timeout", 5*time.Second)
	}

	if readTimeout := os.Getenv("AM_SERVER_READ_TIMEOUT"); readTimeout != "" {
		timeout, err := time.ParseDuration(readTimeout)
		if err != nil {
			logger.Info("you've passed incorrect value of env variable 'AM_SERVER_READ_TIMEOUT', so it will be with default value 5s")
			viper.SetDefault("archive_manager.read_timeout", 5*time.Second)
		} else {
			viper.SetDefault("archive_manager.read_timeout", timeout)
		}
	} else {
		viper.SetDefault("archive_manager.read_timeout", 5*time.Second)
	}

	if idleTimeout := os.Getenv("AM_SERVER_IDLE_TIMEOUT"); idleTimeout != "" {
		timeout, err := time.ParseDuration(idleTimeout)
		if err != nil {
			logger.Info("you've passed incorrect value of env variable 'AM_SERVER_IDLE_TIMEOUT', so it will be with default value 3s")
			viper.SetDefault("archive_manager.idle_timeout", 3*time.Second)
		} else {
			viper.SetDefault("archive_manager.idle_timeout", timeout)
		}
	} else {
		viper.SetDefault("archive_manager.idle_timeout", 3*time.Second)
	}

	if shutdownDuration := os.Getenv("AM_SERVER_SHUTDOWN_DURATION"); shutdownDuration != "" {
		duration, err := time.ParseDuration(shutdownDuration)
		if err != nil {
			logger.Info("you've passed incorrect value of env variable 'AM_SERVER_SHUTDOWN_DURATION', so it will be with default value 10s")
			viper.SetDefault("archive_manager.shutdown_duration", 10*time.Second)
		} else {
			viper.SetDefault("archive_manager.shutdown_duration", duration)
		}
	} else {
		viper.SetDefault("archive_manager.shutdown_duration", 10*time.Second)
	}
}

// Read получает переменные из среды и файла конфигурации
func Read(configFilePath string, logger *zap.Logger) {
	readEnvAndSetDefault(logger)
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
