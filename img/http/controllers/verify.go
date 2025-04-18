package controllers

import (
	"fmt"

	"github.com/IbraheemHaseeb7/fyp-backend/utils"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/labstack/echo/v4"
)

func VerifyCard(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {

		registrationNumber := c.Get("auth_user_registration_number").(string)

		// validate request
		type Request struct {
			File string `json:"file" validate:"required"`
		}
		var requestBody Request
		if err := cr.BindAndValidate(&requestBody, &c); err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		// store in database that the student card is verified
		uuid := watermill.NewUUID()
		utils.Requests[uuid] = make(chan pubsub.PubsubMessage)

		// publishing a read message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "files",
			Operation: "VERIFY_CARD",
			Topic:     "img->db",
			UUID:      uuid,
			Payload: fmt.Sprintf(`{
					"registrationNumber": "%s",
					"cardStatus": "%s"
				}`, registrationNumber, "true"),
		}
		err := cr.Publisher.PublishMessage(pubsubMessage)
		if err != nil {
			return cr.SendErrorResponse(&c)
		}

		cr.GetAndFormResponse(pubsubMessage)

		return cr.SendResponse(&c)
	}
}

func VerifySelfie(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {

		registrationNumber := c.Get("auth_user_registration_number").(string)

		// validate request
		type Request struct {
			File string `json:"file" validate:"required"`
		}
		var requestBody Request
		if err := cr.BindAndValidate(&requestBody, &c); err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		// store in database that the student card is verified
		uuid := watermill.NewUUID()
		utils.Requests[uuid] = make(chan pubsub.PubsubMessage)

		// publishing a read message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "files",
			Operation: "VERIFY_SELFIE",
			Topic:     "img->db",
			UUID:      uuid,
			Payload: fmt.Sprintf(`{
					"registrationNumber": "%s",
					"cardStatus": "%s"
				}`, registrationNumber, "true"),
		}
		err := cr.Publisher.PublishMessage(pubsubMessage)
		if err != nil {
			return cr.SendErrorResponse(&c)
		}

		cr.GetAndFormResponse(pubsubMessage)

		return cr.SendResponse(&c)
	}
}
