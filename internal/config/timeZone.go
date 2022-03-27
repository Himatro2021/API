package config

import (
	"himatro-api/internal/util"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func TimeZone() string {
	tz := os.Getenv("TZ")

	if tz == "" {
		util.LogErr("WARN", "TZ is not found in the env", "")
		log.Println("Unable to locate timezone config in .env, using default value...")

		return "Asia/Jakarta"
	}

	return tz
}
