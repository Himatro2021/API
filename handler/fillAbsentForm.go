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

	updateToken, err := controller.FillAbsentForm(absentID, NPM, keterangan)

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: fmt.Sprintf("Failed to fill absent form because: %s", err.Error()),
		})
	}

	updateTokenExpiresSec, err := strconv.Atoi(os.Getenv("UPDATE_ABSENT_LIST_TOKEN_EXP_SEC"))

	if err != nil {
		updateTokenExpiresSec = 3600 // 1 hour expiry
	}

	cookie := new(http.Cookie)
	cookie.Name = os.Getenv("UPDATE_ABSENT_LIST_COOKIE_NAME")
	cookie.Value = updateToken
	cookie.Expires = time.Now().Add(time.Second * time.Duration(updateTokenExpiresSec))

	c.SetCookie(cookie)

	return nil
}

func UpdateAbsentListByAttendant(c echo.Context) error {
	cookie, err := c.Cookie(os.Getenv("UPDATE_ABSENT_LIST_COOKIE_NAME"))

	if err != nil {
		return c.JSON(http.StatusForbidden, ErrorMessage{
			OK:      false,
			Message: "Please provide update absent token.",
		})
	}

	keterangan := c.FormValue("keterangan")
	absentID, err := strconv.Atoi(c.Param("absentID"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: "absentID must be a valid numeric string",
		})
	}

	if keterangan == "" || (keterangan != "h" && keterangan != "i") {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: "All required field must not empty and use only valid value.",
		})
	}

	if err := controller.UpdateAbsentListByAttendant(absentID, keterangan, cookie); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: fmt.Sprintf("Failed to update absent list because: %s", err.Error()),
		})
	}

	return c.NoContent(http.StatusAccepted)
}
