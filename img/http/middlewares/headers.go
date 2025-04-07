package middlewares

import "github.com/labstack/echo/v4"

func Headers() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			c.Response().Header().Add("Content-Type", "application/json")

			return next(c)
		}
	}
}
