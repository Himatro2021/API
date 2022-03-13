package handler

import (
	"fmt"
	"himatro-api/controller"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetAbsentResult(c echo.Context) error {
	absentID, err := strconv.Atoi(c.Param("absentID"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: "Param: absentID must be a valid numeric string.",
		})
	}

	absentList, err := controller.GetAbsentListResult(absentID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: fmt.Sprintf("Failed to get requested absent list because: %s", err.Error()),
		})
	}

	return c.JSON(http.StatusOK, SuccessListAbsent{
		OK:     true,
		FormID: absentID,
		Total:  len(absentList),
		List:   absentList,
	})
}
