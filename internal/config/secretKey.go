package config

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var defaultSecretKey = "THISISTHEONLYDEFAULTKEYTHATSHOULDNEVERBEUSEDPLEASEREMEMBERTHAT"

func SecretKey() string {
	secretKey := os.Getenv("SECRET_KEY")

	if secretKey == "" {
		log.Println("WARNING! No SECRET_KEY provided. Using the default one is not save at all.")

		return defaultSecretKey
	}

	return secretKey
}
