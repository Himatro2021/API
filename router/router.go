package router

import (
	"himatro-api/handler"
	"himatro-api/middleware"

	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
)

func Router() *echo.Echo {
	e := echo.New()

	e.Use(echoMiddleware.Logger())

	e.GET("/", handler.HomeGet)
	e.GET("/absensi/:absentID", handler.GetAbsentList)
	e.POST("/login", handler.Login)

	e.GET("/admin", handler.Admin)
	e.POST("/admin/absensi", handler.InitAbsent, middleware.RequireLogin)

	return e
}
