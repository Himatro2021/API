package handler

import (
	"himatro-api/internal/contract"
	"himatro-api/internal/controller"
	"himatro-api/internal/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

func InitAbsent(c echo.Context) error {
	createAbsentPayload := contract.CreateAbsentForm{}

	if err := c.Bind(&createAbsentPayload); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: "Invalid type of JSON payload received",
		})
	}

	if err := util.Validator.Struct(&createAbsentPayload); err != nil {
		return c.JSON(http.StatusOK, JSONPayloadValidationError{
			OK:      false,
			Message: "JSON payload validation error",
			Details: util.ExtractValidationErrorMsg(err),
		})
	}

	initAbsentPayload, err := controller.ExtractInitAbsentPayload(createAbsentPayload)

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: err.Error(),
		})
	}

	absentID, err := controller.RegisterNewAbsentForm(&initAbsentPayload)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorMessage{
			OK:      false,
			Message: err.Error(),
		})
	}

	err = controller.InitAbsentList(&initAbsentPayload, absentID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorMessage{
			OK:      false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, SuccessCreateAbsent{
		OK:                          true,
		AbsentID:                    absentID,
		Title:                       initAbsentPayload.Title,
		Participant:                 initAbsentPayload.Participant,
		StartAt:                     initAbsentPayload.StartAt,
		FinishAt:                    initAbsentPayload.FinishAt,
		RequireAttendanceImageProof: initAbsentPayload.RequireAttendanceImageProof,
		RequireExecuseImageProof:    initAbsentPayload.RequireExecuseImageProof,
	})
}
