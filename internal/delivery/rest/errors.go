package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	// ErrBadRequest self explained
	ErrBadRequest = echo.NewHTTPError(http.StatusBadRequest, "Invalid request. Please send a valid payload")

	// ErrValidation self explained
	ErrValidation = echo.NewHTTPError(http.StatusBadRequest, "Validation Error")

	// ErrNotFound self explained
	ErrNotFound = echo.NewHTTPError(http.StatusNotFound, "Record Not Found")

	// ErrInternal self explained
	ErrInternal = echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")

	// ErrForbidden self explained
	ErrForbidden = echo.NewHTTPError(http.StatusForbidden, "Access Forbidden")

	// ErrAlreadyExists self explained
	ErrAlreadyExists = echo.NewHTTPError(http.StatusForbidden, "Record already exists!")
)
