package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HomeGet(c echo.Context) error {
	return c.String(http.StatusOK, "hello world")
}
