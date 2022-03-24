package db

import (
	"himatro-api/internal/config"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func Connect() {
	dsn := config.DBConnString()

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to Postgres Server")
		panic("Connection to Postgres Server failed")
	}

	log.Print("Successfully connected to Postgres Server")
}
