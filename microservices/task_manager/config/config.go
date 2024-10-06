package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// readEnvAndSetDefault устанавливает переменные конфигурации viper по умолчанию. Используется для случая,
// когда файл конфигурации не был найден. Использует переменные окружения для настройки.
func readEnvAndSetDefault(logger *zap.Logger) {
	// CLIENTS HOSTS AND PORTS
	viper.SetDefault("microservice_privelege.host", os.Getenv("PS_SERVER_CONNECTION_HOST"))
	viper.SetDefault("microservice_privelege.port", os.Getenv("PS_SERVER_PORT"))

	viper.SetDefault("microservice_archive.host", os.Getenv("AM_SERVER_CONNECTION_HOST"))
	viper.SetDefault("microservice_archive.port", os.Getenv("AM_SERVER_PORT"))
	// SERVER
	viper.SetDefault("task_manager.address", os.Getenv("TM_SERVER_ADDRESS"))
	if writeTimeout := os.Getenv("TM_SERVER_WRITE_TIMEOUT"); writeTimeout != "" {
		timeout, err := time.ParseDuration(writeTimeout)
		if err != nil {
			logger.Info("you've passed incorrect value of env variable 'TM_SERVER_WRITE_TIMEOUT', so it will be with default value 5s")
			viper.SetDefault("task_manager.write_timeout", 5*time.Second)
		} else {
			viper.SetDefault("task_manager.write_timeout", timeout)
		}
	} else {
		viper.SetDefault("task_manager.write_timeout", 5*time.Second)
	}

	if readTimeout := os.Getenv("TM_SERVER_READ_TIMEOUT"); readTimeout != "" {
		timeout, err := time.ParseDuration(readTimeout)
		if err != nil {
			logger.Info("you've passed incorrect value of env variable 'TM_SERVER_READ_TIMEOUT', so it will be with default value 5s")
			viper.SetDefault("task_manager.read_timeout", 5*time.Second)
		} else {
			viper.SetDefault("task_manager.read_timeout", timeout)
		}
	} else {
		viper.SetDefault("task_manager.read_timeout", 5*time.Second)
	}

	if idleTimeout := os.Getenv("TM_SERVER_IDLE_TIMEOUT"); idleTimeout != "" {
		timeout, err := time.ParseDuration(idleTimeout)
		if err != nil {
			logger.Info("you've passed incorrect value of env variable 'TM_SERVER_IDLE_TIMEOUT', so it will be with default value 3s")
			viper.SetDefault("task_manager.idle_timeout", 3*time.Second)
		} else {
			viper.SetDefault("task_manager.idle_timeout", timeout)
		}
	} else {
		viper.SetDefault("task_manager.idle_timeout", 3*time.Second)
	}

	if shutdownDuration := os.Getenv("TM_SERVER_SHUTDOWN_DURATION"); shutdownDuration != "" {
		duration, err := time.ParseDuration(shutdownDuration)
		if err != nil {
			logger.Info("you've passed incorrect value of env variable 'TM_SERVER_SHUTDOWN_DURATION', so it will be with default value 10s")
			viper.SetDefault("task_manager.shutdown_duration", 10*time.Second)
		} else {
			viper.SetDefault("task_manager.shutdown_duration", duration)
		}
	} else {
		viper.SetDefault("task_manager.shutdown_duration", 10*time.Second)
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
