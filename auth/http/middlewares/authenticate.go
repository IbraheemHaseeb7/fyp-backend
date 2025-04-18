package middlewares

import (
	"net/http"
	"strings"

	"github.com/IbraheemHaseeb7/fyp-backend/auth"
	"github.com/labstack/echo/v4"
)

func Authenticate() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authArray := strings.Split(c.Request().Header.Get("Authorization"), " ")

			if len(authArray) < 2 {
				return c.JSON(http.StatusUnauthorized, map[string]any{
					"data":   nil,
					"status": "Not able to extract token",
					"error":  "Auth token not found in the auth headers",
				})
			}

			token := authArray[1]

			if err := auth.VerifyToken(token); err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]any{
					"data":   nil,
					"status": "Token invalid",
					"error":  err.Error(),
				})
			}

			claims, err := auth.GetClaimsFromToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]any{
					"data":   nil,
					"status": "Not able to extract claims",
					"error":  err.Error(),
				})
			}

			c.Set("auth_user_id", claims["id"])
			c.Set("auth_user_email", claims["email"])
			c.Set("auth_user_name", claims["name"])
			c.Set("auth_user_registration_number", claims["registrationNumber"])
			return next(c)
		}
	}
}
