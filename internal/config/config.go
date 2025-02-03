package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	// API credentials
	APIID    int
	APIHash  string
	BotToken string

	// File paths
	SessionFile string

	// Default settings
	DefaultSpyUserID int64

	// Logging
	LogLevel string
)

// LoadConfig загружает конфигурацию из .env файла
func LoadConfig() error {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Printf("Предупреждение: файл .env не найден: %v", err)
	}

	var err error

	// API credentials
	apiIDStr := os.Getenv("TELEGRAM_API_ID")
	APIID, err = strconv.Atoi(apiIDStr)
	if err != nil {
		return err
	}

	APIHash = os.Getenv("TELEGRAM_API_HASH")
	if APIHash == "" {
		return ErrMissingEnvVar("TELEGRAM_API_HASH")
	}

	BotToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	if BotToken == "" {
		return ErrMissingEnvVar("TELEGRAM_BOT_TOKEN")
	}

	// User settings
	spyUserIDStr := os.Getenv("DEFAULT_SPY_USER_ID")
	DefaultSpyUserID, err = strconv.ParseInt(spyUserIDStr, 10, 64)
	if err != nil {
		return err
	}

	// File paths
	SessionFile = os.Getenv("SESSION_FILE")
	if SessionFile == "" {
		SessionFile = "session.data" // значение по умолчанию
	}

	// Logging
	LogLevel = os.Getenv("LOG_LEVEL")
	if LogLevel == "" {
		LogLevel = "info" // значение по умолчанию
	}

	return nil
}

// ErrMissingEnvVar возвращает ошибку о отсутствующей переменной окружения
func ErrMissingEnvVar(name string) error {
	return &MissingEnvVarError{VarName: name}
}

// MissingEnvVarError представляет ошибку отсутствующей переменной окружения
type MissingEnvVarError struct {
	VarName string
}

func (e *MissingEnvVarError) Error() string {
	return "отсутствует обязательная переменная окружения: " + e.VarName
}