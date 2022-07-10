package config

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// PostgresDSN get DSN string for postgres connection
func PostgresDSN() string {
	host := os.Getenv("PG_HOST")
	db := os.Getenv("PG_DATABASE")
	user := os.Getenv("PG_USER")
	pw := os.Getenv("PG_PASSWORD")
	port := os.Getenv("PG_PORT")

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, pw, db, port)
}

// LogLevel get log level for all server logger instances
func LogLevel() string {
	return os.Getenv("LOG_LEVEL")
}

// ServerPort get desired http port to be running on
func ServerPort() string {
	cfg := os.Getenv("SERVER_PORT")
	if cfg == "" {
		logrus.Warn("Failed to lookup SERVER_PORT env. using default value")
		return "5000" // default port
	}

	return cfg
}
