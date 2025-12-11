package repository

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "crmstom_user"),
		Password: getEnv("DB_PASSWORD", "crmstom_password"),
		DBName:   getEnv("DB_NAME", "crmstom"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func (config *DatabaseConfig) GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)
}

func ConnectToDatabase(config *DatabaseConfig) (*sql.DB, error) {
	connStr := config.GetConnectionString()

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	// Проверка подключения
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка ping базы данных: %w", err)
	}

	return db, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
