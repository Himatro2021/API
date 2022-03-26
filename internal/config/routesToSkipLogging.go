package config

import (
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

func GetSkippedLogRoutes() []string {
	raw := os.Getenv("LOG_SKIPPED_ROUTES")

	if raw == "" {
		return []string{"/login"} // default skipped routes log
	}

	return strings.Split(raw, ",")
}
