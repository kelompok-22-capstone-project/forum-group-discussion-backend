package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func BodyLimit(e *echo.Echo) {
	e.Use(middleware.BodyLimit("128K"))
}
