package middlewares

import (
	// "strings"
	"github.com/labstack/echo/v4"
)

func Authenticate() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// token := strings.Split(c.Request().Header.Get("Authorization"), " ")[1]

			return next(c)
		}
	}
}
