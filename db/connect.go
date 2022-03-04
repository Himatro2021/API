package db

import (
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func Connect() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", os.Getenv("PG_HOST"), os.Getenv("PG_USER"), os.Getenv("PG_PASSWORD"), os.Getenv("PG_DATABASE"), os.Getenv("PG_PORT"))

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to Postgres Server")
		panic("Connection to Postgres Server failed")
	}

	log.Print("Successfully connected to Postgres Server")
}
