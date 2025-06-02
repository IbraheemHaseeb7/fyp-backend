package controllers

import (
	"encoding/json"

	"github.com/IbraheemHaseeb7/fyp-backend/utils"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/labstack/echo/v4"
)

func GetRoom(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		roomID := c.Param("id")
		if roomID == "" {
			return c.JSON(400, map[string]string{"error": "room ID is required"})
		}

		uuid := watermill.NewUUID()
		utils.Requests.Store(uuid, make(chan pubsub.PubsubMessage))

		payload, err := json.Marshal(map[string]any{"room_id": roomID}); if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		// publishing a read message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "rooms",
			Operation: "READ_ONE",
			Topic:     "auth->db",
			UUID:      uuid,
			Payload:   string(payload),
		}
		err = cr.Publisher.PublishMessage(pubsubMessage)
		if err != nil {
			return cr.SendErrorResponse(&c)
		}

		cr.GetAndFormResponse(pubsubMessage)
		return cr.SendResponse(&c)
	}
}

func GetRoomWithIDs(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestID := c.QueryParam("request_id")
		proposalID := c.QueryParam("proposal_id")

		if requestID == "" || proposalID == "" {
			return c.JSON(400, map[string]string{"error": "request_id and proposal_id are required"})
		}

		uuid := watermill.NewUUID()
		utils.Requests.Store(uuid, make(chan pubsub.PubsubMessage))

		payload, err := json.Marshal(map[string]any{"request_id": requestID, "proposal_id": proposalID}); if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		// publishing a read message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "rooms",
			Operation: "READ_WITH_IDS",
			Topic:     "auth->db",
			UUID:      uuid,
			Payload:   string(payload),
		}
		err = cr.Publisher.PublishMessage(pubsubMessage)
		if err != nil {
			return cr.SendErrorResponse(&c)
		}

		cr.GetAndFormResponse(pubsubMessage)
		return cr.SendResponse(&c)
	}
}
