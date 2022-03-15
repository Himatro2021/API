package handler

import (
	"fmt"
	"himatro-api/controller"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
)

func CheckAbsentForm(c echo.Context) error {
	absentID, err := strconv.Atoi(c.Param("absentID"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: "Absent ID must be a valid numeric string.",
		})
	}

	if err := controller.IsFormWriteable(absentID); err != nil {
		return c.JSON(http.StatusForbidden, ErrorMessage{
			OK:      false,
			Message: fmt.Sprintf("Request forbidden because: %s", err.Error()),
		})
	}

	return c.NoContent(http.StatusOK)
}

func FillAbsentForm(c echo.Context) error {
	absentID, err := strconv.Atoi(c.Param("absentID"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: "Absent ID must be a valid numeric string.",
		})
	}

	if err := controller.IsFormWriteable(absentID); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: fmt.Sprintf("Failed to fill absent form because: %s", err.Error()),
		})
	}

	NPM := c.FormValue("NPM")
	keterangan := c.FormValue("keterangan")

	if NPM == "" || keterangan == "" {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: "All required field must be supplied and using valid value.",
		})
	}

	if keterangan != "h" && keterangan != "i" {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: fmt.Sprintf("Value keterangan: %s is invalid.", keterangan),
		})
	}

	if err := controller.FillAbsentForm(absentID, NPM, keterangan); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: fmt.Sprintf("Failed to fill absent form because: %s", err.Error()),
		})
	}

	return nil
}
