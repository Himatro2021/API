package handler

import (
	"fmt"
	"himatro-api/db"
	"himatro-api/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ListAbsensi(c echo.Context) error {
	var res []models.AnggotaBiasa

	query := db.DB.Where("npm != ?", "1").Find(&res)

	if query.Error != nil {
		response := ErrorMessage{
			OK:      false,
			Message: "Query error.",
		}

		return c.JSON(http.StatusBadRequest, response)
	}

	fmt.Println()

	response := AbsentListSuccessMessage{
		OK:     true,
		Status: 200,
		Result: res,
	}

	return c.JSON(http.StatusOK, response)
}
