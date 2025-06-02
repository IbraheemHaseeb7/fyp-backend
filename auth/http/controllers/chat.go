package controllers

import (
	"encoding/json"

	"github.com/IbraheemHaseeb7/fyp-backend/utils"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/labstack/echo/v4"
)

func ChatMessages(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Extract the chat ID from the request parameters
		roomID := c.Param("room_id")
		if roomID == "" {
			return c.JSON(400, map[string]string{"error": "chat_id is required"})
		}

		pageNo := c.QueryParam("page")
		if pageNo == "" {
			pageNo = "1" // Default to page 1 if not provided
		}

		uuid := watermill.NewUUID()
		utils.Requests.Store(uuid, make(chan pubsub.PubsubMessage))

		// Prepare the payload with chat ID and page number
		payload, err := json.Marshal(map[string]any{
			"room_id": roomID,
			"page":    pageNo,
			"user_id": c.Get("auth_user_id"),
		}); if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		// publishing a read message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "chats",
			Operation: "GET_MESSAGES",
			Topic:     "auth->db",
			UUID:      uuid,
			Payload:   string(payload),
		}
		err = cr.Publisher.PublishMessage(pubsubMessage); if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		cr.GetAndFormResponse(pubsubMessage)
		return cr.SendResponse(&c)

	}
}
