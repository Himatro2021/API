package config

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	// DefaultTokenLength ruled how long the token should be generated
	DefaultTokenLength = 90

	// DefaultAccessTokenExpiry ruled how long the access token expired in hour
	DefaultAccessTokenExpiry = time.Hour * 24 * 7

	// DefaultRefreshTokenExpiry ruled how long the refresh token expired in hour
	DefaultRefreshTokenExpiry = time.Hour * 24 * 30
)

// Env get env value for ENV key
func Env() string {
	cfg := os.Getenv("ENV")

	if cfg == "" {
		return "development"
	}

	return cfg
}

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

// PrivateKey get private key from env
func PrivateKey() string {
	key := os.Getenv("PRIVATE_KEY")
	if key == "" {
		logrus.Error("PRIVATE_KEY is unset. May cause danger in encryption method")
	}

	return key
}

// IvKey get private key from env
func IvKey() string {
	key := os.Getenv("IV_KEY")
	if key == "" {
		logrus.Error("IV_KEY is unset. May cause danger in encryption method")
	}

	return key
}
