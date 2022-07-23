package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	// ErrNotFound represent condition when auth process encounter error not found
	ErrNotFound = echo.NewHTTPError(http.StatusNotFound, "Not Found!")

	// ErrInternal internal error when performing auth process
	ErrInternal = echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error.")

	// ErrUnauthorized most common representation when auth process decide a failure in the auth process
	ErrUnauthorized = echo.NewHTTPError(http.StatusUnauthorized, "Request Unauthorized!")
)
