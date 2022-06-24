package middleware

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func JWTMiddleware() echo.MiddlewareFunc {
	secret := os.Getenv("JWT_SECRET")
	config := middleware.JWTConfig{
		SigningKey: []byte(secret),
	}

	return middleware.JWTWithConfig(config)
}
