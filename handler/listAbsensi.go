package handler

import (
	"himatro-api/controller"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetAbsentList(c echo.Context) error {
	absentID, err := strconv.Atoi(c.Param("absentID"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: "Invalid absentID. Must be a number.",
		})
	}

	absentList, err := controller.GetAbsentList(absentID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorMessage{
			OK:      false,
			Message: err.Error(),
		})
	}

	if len(absentList) == 0 {
		return c.JSON(http.StatusNotFound, ErrorMessage{
			OK:      false,
			Message: "Absent list is not found. Please use valid absentID.",
		})
	}

	return c.JSON(http.StatusOK, SuccessListAbsent{
		OK:     true,
		FormID: absentID,
		Total:  len(absentList),
		List:   absentList,
	})
}
