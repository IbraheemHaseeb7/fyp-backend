package middlewares

import (
	"strings"

	"github.com/IbraheemHaseeb7/fyp-backend/http/controllers"
	"github.com/IbraheemHaseeb7/fyp-backend/utils"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/labstack/echo/v4"
)

func Authenticate(cr controllers.ControllerRequest) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			tokenHeader := strings.Split(c.Request().Header.Get("Authorization"), " ")
			if len(tokenHeader) > 1 {
				uuid := watermill.NewUUID()
				utils.Requests[uuid] = make(chan pubsub.PubsubMessage)

				cr.Publisher.PublishMessage(pubsub.PubsubMessage{
					Payload: map[string]string{
						"token": tokenHeader[1],
					},
					Entity:    "files",
					Operation: "GET_CLAIMS",
					UUID:      uuid,
					Topic:     "img->auth",
				})
				authResp := <-utils.Requests[uuid]
				delete(utils.Requests, uuid)

				payload, ok := authResp.Payload.(map[string]any)
				if !ok {
					return cr.SendErrorResponse(&c)
				}

				if status := payload["verified"]; status == false {
					cr.APIResponse.Status = "Token invalid"
					cr.APIResponse.StatusCode = 401
					return cr.SendErrorResponse(&c)
				}

				c.Set("auth_user_id", payload["id"])
				c.Set("auth_user_email", payload["email"])
				c.Set("auth_user_name", payload["name"])
				c.Set("auth_user_registration_number", payload["registrationNumber"])

				return next(c)
			}

			cr.APIResponse.Status = "Token not found"
			cr.APIResponse.StatusCode = 401
			return cr.SendErrorResponse(&c)
		}
	}
}
