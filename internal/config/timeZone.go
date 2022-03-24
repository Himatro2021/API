package config

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func TimeZone() string {
	tz := os.Getenv("TZ")

	if tz == "" {
		log.Println("Unable to locate timezone config in .env, using default value...")

		return "Asia/Jakarta"
	}

	return tz
}
