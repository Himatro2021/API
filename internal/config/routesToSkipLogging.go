package config

import (
	"himatro-api/internal/util"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

func GetSkippedLogRoutes() []string {
	raw := os.Getenv("LOG_SKIPPED_ROUTES")

	if raw == "" {
		util.LogErr("WARN", "LOG_SKIPPED_ROUTES is not found on env", "")
		return []string{"/login"} // default skipped routes log
	}

	return strings.Split(raw, ",")
}
