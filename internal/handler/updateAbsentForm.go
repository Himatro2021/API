package handler

import (
	"fmt"
	"himatro-api/internal/contract"
	"himatro-api/internal/controller"
	"himatro-api/internal/util"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func UpdateFormTitle(c echo.Context) error {
	payload := contract.UpdateFormTitle{}

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

	absentID, err := strconv.Atoi(c.Param("absentID"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: fmt.Sprintf("AbsentID with ID %s is invalid.", c.Param("absentID")),
		})
	}

	newValue, err := controller.UpdateFormTitle(absentID, payload.Title)

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

	payload := contract.UpdateFormParticipant{}

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

	if err = controller.UpdateParticipant(absentID, payload.Participant); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: fmt.Sprintf("Update participant in absent form is failed because: %s", err.Error()),
		})
	}

	return c.JSON(http.StatusOK, SuccessUpdateForm{
		OK:        true,
		Message:   "Update success",
		FieldName: "participant",
		Value:     strings.ToUpper(payload.Participant),
	})
}

func UpdateFormStartAt(c echo.Context) error {
	absentID, err := strconv.Atoi(c.Param("absentID"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: "AbsentID must be a valid number.",
		})
	}

	payload := contract.UpdateFormTime{}

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

	startTime, err := controller.UpdateAbsentFormStartAt(absentID, payload.Date, payload.Time)

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
	absentID, err := strconv.Atoi(c.Param("absentID"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: "AbsentID must be a valid number.",
		})
	}

	payload := contract.UpdateFormTime{}

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

	finishTime, err := controller.UpdateAbsentFormFinishAt(absentID, payload.Date, payload.Time)

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

	payload := contract.UpdateFormImageProof{}

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

	if err := controller.UpdateAbsentFormExecuseImageProof(absentID, payload.Status); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: fmt.Sprintf("Failed to update attendance image proof for Absent Form with ID: %d because: %s", absentID, err.Error()),
		})
	}

	return c.JSON(http.StatusOK, SuccessUpdateForm{
		OK:        true,
		Message:   "Update Success",
		FieldName: "requireExecuseImageProof",
		Value:     fmt.Sprintf("%t", payload.Status),
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

	payload := contract.UpdateFormImageProof{}

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

	if err := controller.UpdateAbsentFormAttendanceImageProof(absentID, payload.Status); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: fmt.Sprintf("Failed to update attendance image proof for Absent Form with ID: %d because: %s", absentID, err.Error()),
		})
	}

	return c.JSON(http.StatusOK, SuccessUpdateForm{
		OK:        true,
		Message:   "Update Success",
		FieldName: "requireAttendanceImageProof",
		Value:     fmt.Sprintf("%t", payload.Status),
	})
}
