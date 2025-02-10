package controllers

import (
	"github.com/IbraheemHaseeb7/fyp-backend/utils"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/labstack/echo/v4"
)

func CreateUsers(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")

		uuid := watermill.NewUUID()
		utils.Requests[uuid] = make(chan pubsub.PubsubMessage)

		// publishing a create message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "users",
			Operation: "CREATE",
			Topic:     "auth->db",
			UUID:      uuid,
			Payload: `
			{
				"name": "Ibraheem",
				"email": "ibraheemhaseeb7@gmail.com",
				"registrationNumber": "FA21-BCS-052",
				"dob": "2003-04-07",
				"semester": 7,
				"department": "CS"
			}`,
		}
		err := cr.Publisher.PublishMessage(pubsubMessage)
		if err != nil {
			return cr.SendErrorResponse(&c)
		}

		cr.GetAndFormResponse(pubsubMessage)
		return cr.SendResponse(&c)
	}
}

func ReadAllUsers(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")

		uuid := watermill.NewUUID()
		utils.Requests[uuid] = make(chan pubsub.PubsubMessage)

		// publishing a read message
		pubsubsMessage := pubsub.PubsubMessage{
			Entity:    "users",
			Operation: "READ_ALL",
			Topic:     "auth->db",
			UUID:      uuid,
		}
		err := cr.Publisher.PublishMessage(pubsubsMessage)
		if err != nil {
			return cr.SendErrorResponse(&c)
		}

		cr.GetAndFormResponse(pubsubsMessage)
		return cr.SendResponse(&c)
	}
}

func ReadOneUser(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		regNo := c.Param("id")

		uuid := watermill.NewUUID()
		utils.Requests[uuid] = make(chan pubsub.PubsubMessage)

		// publishing a read message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "users",
			Operation: "READ_ONE",
			Topic:     "auth->db",
			UUID:      uuid,
			Payload:   `{"registrationNumber": "` + regNo + `"}`,
		}
		err := cr.Publisher.PublishMessage(pubsubMessage)
		if err != nil {
			return cr.SendErrorResponse(&c)
		}

		cr.GetAndFormResponse(pubsubMessage)
		return cr.SendResponse(&c)
	}
}
