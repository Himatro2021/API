package handler

import (
	"fmt"
	"himatro-api/internal/config"
	"himatro-api/internal/contract"
	"himatro-api/internal/controller"
	"himatro-api/internal/util"
	"net/http"
	"strconv"
	"time"

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

	payload := contract.FillAbsentList{}

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: "Invalid type of JSON Payload received",
		})
	}

	if err := util.Validator.Struct(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, JSONPayloadValidationError{
			OK:      false,
			Message: "JSON payload validation error",
			Details: util.ExtractValidationErrorMsg(err),
		})
	}

	updateToken, err := controller.FillAbsentForm(absentID, payload.NPM, payload.Keterangan)

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: fmt.Sprintf("Failed to fill absent form because: %s", err.Error()),
		})
	}

	updateTokenExpiresSec := config.UpdateAbsentListTokenExpSec()

	cookie := new(http.Cookie)
	cookie.Name = config.UpdateAbsentListCookieName()
	cookie.Value = updateToken
	cookie.Expires = time.Now().Add(time.Second * time.Duration(updateTokenExpiresSec))

	c.SetCookie(cookie)

	return nil
}

func UpdateAbsentListByAttendant(c echo.Context) error {
	cookie, err := c.Cookie(config.UpdateAbsentListCookieName())

	if err != nil {
		return c.JSON(http.StatusForbidden, ErrorMessage{
			OK:      false,
			Message: "Please provide update absent token.",
		})
	}

	absentID, err := strconv.Atoi(c.Param("absentID"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: "absentID must be a valid numeric string",
		})
	}

	payload := contract.UpdateKeteranganAbsent{}

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: "Invalid type of JSON Payload received",
		})
	}

	if err := util.Validator.Struct(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, JSONPayloadValidationError{
			OK:      false,
			Message: "JSON payload validation error",
			Details: util.ExtractValidationErrorMsg(err),
		})
	}

	if err := controller.UpdateAbsentListByAttendant(absentID, payload.Keterangan, cookie); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: fmt.Sprintf("Failed to update absent list because: %s", err.Error()),
		})
	}

	return c.NoContent(http.StatusAccepted)
}
