package middleware

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func KeyAuth() echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(
		middleware.KeyAuthConfig{
			KeyLookup: "header:api-key",
			Validator: func(auth string, c echo.Context) (bool, error) { return auth == os.Getenv("API_KEY"), nil },
		},
	)
}
