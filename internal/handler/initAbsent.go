package handler

import (
	"fmt"
	"himatro-api/internal/controller"
	"net/http"

	"github.com/labstack/echo/v4"
)

func InitAbsent(c echo.Context) error {
	initAbsentPayload, err := controller.ExtractInitAbsentPayload(c)

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
		fmt.Println(err)
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
