package middleware

import (
	"himatro-api/internal/config"

	"github.com/labstack/echo/v4/middleware"
)

var RequireLogin = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningKey: []byte(config.JWTSigningKey()),
})
