package config

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func ModifyDefaultResponses(e *echo.Echo) {
	ModifyDefault404Response(e)
}

func ModifyDefault404Response(e *echo.Echo) {
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Status == http.StatusNotFound {
			c.JSON(http.StatusNotFound, map[string]string{
				"error": "Resource not found",
			})
		} else {
			e.DefaultHTTPErrorHandler(err, c)
		}
	}
}
