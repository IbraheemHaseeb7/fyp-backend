package controllers

import (
	"github.com/IbraheemHaseeb7/fyp-backend/utils"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/labstack/echo/v4"
)

func ReadOneVehicle(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		regNo := c.Param("id")

		uuid := watermill.NewUUID()
		utils.Requests[uuid] = make(chan pubsub.PubsubMessage)

		// publishing a read message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "vehicles",
			Operation: "READ_ONE",
			Topic:     "auth->db",
			UUID:      uuid,
			Payload:   `{"id": ` + regNo + `}`,
		}
		err := cr.Publisher.PublishMessage(pubsubMessage)
		if err != nil {
			return cr.SendErrorResponse(&c)
		}

		cr.GetAndFormResponse(pubsubMessage)
		return cr.SendResponse(&c)
	}
}

func ReadAllVehicles(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")

		uuid := watermill.NewUUID()
		utils.Requests[uuid] = make(chan pubsub.PubsubMessage)

		// publishing a read message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "vehicles",
			Operation: "READ_ALL",
			Topic:     "auth->db",
			UUID:      uuid,
		}
		err := cr.Publisher.PublishMessage(pubsubMessage)
		if err != nil {
			return cr.SendErrorResponse(&c)
		}

		cr.GetAndFormResponse(pubsubMessage)
		return cr.SendResponse(&c)
	}
}

func CreateVehicle(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		// regNo := c.Param("id")

		uuid := watermill.NewUUID()
		utils.Requests[uuid] = make(chan pubsub.PubsubMessage)

		// publishing a read message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "vehicles",
			Operation: "CREATE",
			Topic:     "auth->db",
			UUID:      uuid,
			// Payload:   `{"registrationNumber": "`+regNo+`"}`,
		}
		err := cr.Publisher.PublishMessage(pubsubMessage)
		if err != nil {
			return cr.SendErrorResponse(&c)
		}

		cr.GetAndFormResponse(pubsubMessage)
		return cr.SendResponse(&c)
	}
}
