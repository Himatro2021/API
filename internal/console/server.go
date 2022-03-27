package console

import (
	"fmt"
	"himatro-api/internal/db"
	"himatro-api/internal/router"
	"himatro-api/internal/util"
	"log"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
)

var runServer = &cobra.Command{
	Use:   "server",
	Short: "Start the server",
	Long:  "Use this command to start Himatro API HTTP server",
	Run:   InitServer,
}

func init() {
	RootCmd.AddCommand(runServer)
}

func InitServer(cmd *cobra.Command, args []string) {
	db.Connect()

	r := router.Router()
	s := http.Server{
		Addr:    os.Getenv("SERVER_PORT"),
		Handler: r,
	}

	log.Print(fmt.Sprintf("Server listening on port %s", s.Addr))
	err := s.ListenAndServe()

	if err != nil {
		util.LogErr("ERROR", "SERVER failed to start", err.Error())
		log.Fatal("Server failed to starts.", err)
	}
}
