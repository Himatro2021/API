package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Admin(c echo.Context) error {
	fmt.Println("hey")

	return c.JSON(http.StatusOK, "hey")
}
