package middleware

import (
	"himatro-api/internal/config"
	"os"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

var reqLogName = os.Getenv("LOG_FILE_NAME")
var reqLogFile, _ = os.OpenFile(reqLogName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

func RequestLogger() echo.MiddlewareFunc {
	return echoMiddleware.LoggerWithConfig(echoMiddleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			for _, route := range config.GetSkippedLogRoutes() {
				if route == c.Request().RequestURI {
					return true
				}
			}

			return false
		},
		Output: reqLogFile,
	})
}
