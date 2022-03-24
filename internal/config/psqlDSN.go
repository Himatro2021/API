package config

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func DBConnString() string {
	host := os.Getenv("PG_HOST")
	db := os.Getenv("PG_DATABASE")
	user := os.Getenv("PG_USER")
	pw := os.Getenv("PG_PASSWORD")
	port := os.Getenv("PG_PORT")

	if os.Getenv("ENV") != "local" {
		host = "host.docker.internal"
	}

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, pw, db, port)
}
