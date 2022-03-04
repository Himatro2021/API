package router

import (
	"himatro-api/handler"
	"himatro-api/middleware"

	"github.com/labstack/echo/v4"
)

func Router() *echo.Echo {
	e := echo.New()

	e.GET("/", handler.HomeGet)
	e.GET("/list_absensi", handler.ListAbsensi)
	e.POST("/login", handler.Login)

	e.GET("/admin", handler.Admin, middleware.RequireLogin)

	return e
}
