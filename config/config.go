package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/satori/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	rootEmailDefault     = "root@mail.ru"
	rootPasswordDefault  = "Root1234"
	rootFirstNameDefault = "Root"
	rootLastNameDefault  = "Rootov"
)

// readEnvAndSetDefault устанавливает переменные конфигурации viper по умолчанию. Используется для случая,
// когда файл конфигурации не был найден. Использует переменные окружения для настройки.
func readEnvAndSetDefault(logger *zap.Logger) {
	// PROJECT
	if secretKey := os.Getenv("SECRET_KEY"); secretKey != "" {
		viper.SetDefault("secret_key", secretKey)
	} else {
		viper.SetDefault("secret_key", uuid.NewV4().String())
	}

	// POSTGRES
	if port := os.Getenv("POSTGRES_PORT"); port != "" {
		psqlPort, err := strconv.Atoi(port)
		if err != nil {
			logger.Info("you've passed incorrect value of env variable 'POSTGRES_PORT', so it will be with default value 5432")
			viper.SetDefault("postgres.port", 5432)
		} else {
			viper.SetDefault("postgres.port", psqlPort)
		}
	} else {
		viper.SetDefault("postgres.port", 5432)
	}

	if host := os.Getenv("POSTGRES_CONNECTION_HOST"); host != "" {
		viper.SetDefault("postgres.connectionHost", host)
	} else {
		viper.SetDefault("postgres.connectionHost", "localhost")
	}

	if rootEmail := os.Getenv("ROOT_EMAIL"); rootEmail != "" {
		viper.SetDefault("root_email", rootEmail)
	} else {
		viper.SetDefault("root_email", rootEmailDefault)
	}

	if rootPassword := os.Getenv("ROOT_PASSWORD"); rootPassword != "" {
		viper.SetDefault("root_password", rootPassword)
	} else {
		viper.SetDefault("root_password", rootPasswordDefault)
	}

	if rootFirstName := os.Getenv("ROOT_FIRST_NAME"); rootFirstName != "" {
		viper.SetDefault("root_first_name", rootFirstName)
	} else {
		viper.SetDefault("root_first_name", rootFirstNameDefault)
	}

	if rootLastName := os.Getenv("ROOT_LAST_NAME"); rootLastName != "" {
		viper.SetDefault("root_last_name", rootLastName)
	} else {
		viper.SetDefault("root_last_name", rootLastNameDefault)
	}

	viper.SetDefault("postgres.sslmode", "disable")
	// SERVER
	if address := os.Getenv("SERVER_ADDRESS"); address != "" {
		viper.SetDefault("server.address", address)
	} else {
		viper.SetDefault("server.address", "localhost:8000")
	}

	if writeTimeout := os.Getenv("SERVER_WRITE_TIMEOUT"); writeTimeout != "" {
		timeout, err := time.ParseDuration(writeTimeout)
		if err != nil {
			logger.Info("you've passed incorrect value of env variable 'SERVER_WRITE_TIMEOUT', so it will be with default value 5s")
			viper.SetDefault("server.write_timeout", 5*time.Second)
		} else {
			viper.SetDefault("server.write_timeout", timeout)
		}
	} else {
		viper.SetDefault("server.write_timeout", 5*time.Second)
	}

	if readTimeout := os.Getenv("SERVER_READ_TIMEOUT"); readTimeout != "" {
		timeout, err := time.ParseDuration(readTimeout)
		if err != nil {
			logger.Info("you've passed incorrect value of env variable 'SERVER_READ_TIMEOUT', so it will be with default value 5s")
			viper.SetDefault("server.read_timeout", 5*time.Second)
		} else {
			viper.SetDefault("server.read_timeout", timeout)
		}
	} else {
		viper.SetDefault("server.read_timeout", 5*time.Second)
	}

	if idleTimeout := os.Getenv("SERVER_IDLE_TIMEOUT"); idleTimeout != "" {
		timeout, err := time.ParseDuration(idleTimeout)
		if err != nil {
			logger.Info("you've passed incorrect value of env variable 'SERVER_IDLE_TIMEOUT', so it will be with default value 3s")
			viper.SetDefault("server.idle_timeout", 3*time.Second)
		} else {
			viper.SetDefault("server.idle_timeout", timeout)
		}
	} else {
		viper.SetDefault("server.idle_timeout", 3*time.Second)
	}

	if shutdownDuration := os.Getenv("SERVER_SHUTDOWN_DURATION"); shutdownDuration != "" {
		duration, err := time.ParseDuration(shutdownDuration)
		if err != nil {
			logger.Info("you've passed incorrect value of env variable 'SERVER_SHUTDOWN_DURATION', so it will be with default value 10s")
			viper.SetDefault("server.shutdown_duration", 10*time.Second)
		} else {
			viper.SetDefault("server.shutdown_duration", duration)
		}
	} else {
		viper.SetDefault("server.shutdown_duration", 10*time.Second)
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
