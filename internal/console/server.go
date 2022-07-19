package console

import (
	"fmt"

	"github.com/Himatro2021/API/internal/config"
	"github.com/Himatro2021/API/internal/db"
	"github.com/Himatro2021/API/internal/delivery/rest"
	"github.com/Himatro2021/API/internal/helper"
	"github.com/Himatro2021/API/internal/repository"
	"github.com/Himatro2021/API/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
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

// InitServer initialize HTTP server
func InitServer(cmd *cobra.Command, args []string) {
	db.InitializePostgresConn()
	setupLogger()

	sqlDB, err := db.PostgresDB.DB()
	if err != nil {
		logrus.Fatal("unable to start server. reason: ", err.Error())
	}

	defer helper.WrapCloser(sqlDB.Close)

	userRepo := repository.NewUserRepository(db.PostgresDB)
	userUsecase := usecase.NewUserUsecase(userRepo)

	absentRepo := repository.NewAbsentRepository(db.PostgresDB)
	absentUsecase := usecase.NewAbsentUsecase(absentRepo)

	HTTPServer := echo.New()

	HTTPServer.Pre(middleware.AddTrailingSlash())
	HTTPServer.Use(middleware.Logger())

	RESTGroup := HTTPServer.Group("rest")

	rest.InitService(RESTGroup, userUsecase, absentUsecase)

	if err := HTTPServer.Start(fmt.Sprintf(":%s", config.ServerPort())); err != nil {
		logrus.Fatal("unable to start server. reason: ", err.Error())
	}

	logrus.Info("Server running on port: ", config.ServerPort())
}
