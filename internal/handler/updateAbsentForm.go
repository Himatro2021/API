package handler

import (
	"fmt"
	"himatro-api/internal/controller"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func UpdateFormTitle(c echo.Context) error {
	title := c.FormValue("title")
	absentID, err := strconv.Atoi(c.Param("absentID"))

	if title == "" || err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: fmt.Sprintf("Form with ID: %d is not updated since payloads are invalid", absentID),
		})
	}

	newValue, err := controller.UpdateFormTitle(absentID, title)

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: fmt.Sprintf("Form with ID: %d is not updated because: %s", absentID, err.Error()),
		})
	}

	return c.JSON(http.StatusOK, SuccessUpdateForm{
		OK:        true,
		Message:   "Update success",
		FieldName: "title",
		Value:     newValue,
	})
}

func UpdateFormParticipant(c echo.Context) error {
	absentID, err := strconv.Atoi(c.Param("absentID"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: "Update participant in absent form is failed because absentID is not a number.",
		})
	}

	participant := c.FormValue("participant")

	if err = controller.UpdateParticipant(absentID, participant); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: fmt.Sprintf("Update participant in absent form is failed because: %s", err.Error()),
		})
	}

	return c.JSON(http.StatusOK, SuccessUpdateForm{
		OK:        true,
		Message:   "Update success",
		FieldName: "participant",
		Value:     participant,
	})
}

func UpdateFormStartAt(c echo.Context) error {
	startAtTime := c.FormValue("startAtTime")
	startAtDate := c.FormValue("startAtDate")

	if startAtDate == "" || startAtTime == "" {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: "Failed to update form startAt because missing required payload.",
		})
	}

	absentID, err := strconv.Atoi(c.Param("absentID"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: "AbsentID must be a valid number.",
		})
	}

	startTime, err := controller.UpdateAbsentFormStartAt(absentID, startAtDate, startAtTime)

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: fmt.Sprintf("Update form startAt failed because: %s", err.Error()),
		})
	}

	return c.JSON(http.StatusOK, SuccessUpdateForm{
		OK:        true,
		Message:   "Update Success",
		FieldName: "startAt",
		Value:     startTime,
	})
}

func UpdateAbsentFormFinishAt(c echo.Context) error {
	finishAtTime := c.FormValue("finishAtTime")
	finishAtDate := c.FormValue("finishAtDate")

	if finishAtTime == "" || finishAtDate == "" {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: "Failed to update form finishAt because missing required payload.",
		})
	}

	absentID, err := strconv.Atoi(c.Param("absentID"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: "AbsentID must be a valid number.",
		})
	}

	finishTime, err := controller.UpdateAbsentFormFinishAt(absentID, finishAtDate, finishAtTime)

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: fmt.Sprintf("Update form startAt failed because: %s", err.Error()),
		})
	}

	return c.JSON(http.StatusOK, SuccessUpdateForm{
		OK:        true,
		Message:   "Update Success",
		FieldName: "finishAt",
		Value:     finishTime,
	})
}

func UpdateAbsentFormExecuseImageProof(c echo.Context) error {
	absentID, err := strconv.Atoi(c.Param("absentID"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: "AbsentID must be a valid number.",
		})
	}

	requiredStatus, err := strconv.ParseBool(c.FormValue("status"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: fmt.Sprintf("status : %s is invalid. Please refer to documentation section.", c.FormValue("status")),
		})
	}

	if err := controller.UpdateAbsentFormExecuseImageProof(absentID, requiredStatus); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: fmt.Sprintf("Failed to update attendance image proof for Absent Form with ID: %d because: %s", absentID, err.Error()),
		})
	}

	return c.JSON(http.StatusOK, SuccessUpdateForm{
		OK:        true,
		Message:   "Update Success",
		FieldName: "requireExecuseImageProof",
		Value:     fmt.Sprintf("%t", requiredStatus),
	})
}

func UpdateAbsentFormAttendanceImageProof(c echo.Context) error {
	absentID, err := strconv.Atoi(c.Param("absentID"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: "AbsentID must be a valid number.",
		})
	}

	requiredStatus, err := strconv.ParseBool(c.FormValue("status"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: fmt.Sprintf("status : %s is invalid. Please refer to documentation section.", c.FormValue("status")),
		})
	}

	if err := controller.UpdateAbsentFormAttendanceImageProof(absentID, requiredStatus); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: fmt.Sprintf("Failed to update attendance image proof for Absent Form with ID: %d because: %s", absentID, err.Error()),
		})
	}

	return c.JSON(http.StatusOK, SuccessUpdateForm{
		OK:        true,
		Message:   "Update Success",
		FieldName: "requireAttendanceImageProof",
		Value:     fmt.Sprintf("%t", requiredStatus),
	})
}
