package middleware

import (
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/labstack/echo/v4/middleware"
)

var RequireLogin = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningKey: []byte(os.Getenv("JWT_SECRET_SIGNING_KEY")),
})
