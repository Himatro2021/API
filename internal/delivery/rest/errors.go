package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	// ErrBadRequest self explained
	ErrBadRequest = echo.NewHTTPError(http.StatusBadRequest, "Invalid request. Please send a vaid payload")

	// ErrValidation self explained
	ErrValidation = echo.NewHTTPError(http.StatusBadRequest, "Validation Error")

	// ErrInternal self explained
	ErrInternal = echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
)

func customValidationErrMessage(m string) *echo.HTTPError {
	type validationErr struct {
		Message string `json:"message"`
		Hint    string `json:"hint"`
	}

	return echo.NewHTTPError(http.StatusBadRequest, validationErr{
		Message: "Validation Error",
		Hint:    m,
	})
}
