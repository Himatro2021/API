package server

import (
	"fmt"
	"himatro-api/internal/db"
	"himatro-api/internal/router"
	"log"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func InitServer() {
	db.Connect()

	r := router.Router()
	s := http.Server{
		Addr:    os.Getenv("SERVER_PORT"),
		Handler: r,
	}

	log.Print(fmt.Sprintf("Server listening on port %s", s.Addr))
	err := s.ListenAndServe()

	if err != nil {
		log.Fatal("Server failed to starts.", err)
	}
}
