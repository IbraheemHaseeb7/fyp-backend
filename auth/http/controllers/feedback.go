package controllers

import (
	"encoding/json"
	"time"

	"github.com/IbraheemHaseeb7/fyp-backend/utils"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/labstack/echo/v4"
)

func AddFeedback(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Get("auth_user_id").(float64)

		uuid := watermill.NewUUID()
		utils.Requests.Store(uuid, make(chan pubsub.PubsubMessage))
		var feedback map[string]any
		if err := c.Bind(&feedback); err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}
		feedback["user_id"] = userId
		feedback["created_at"] = time.Now().Format(time.RFC3339)

		payload, err := json.Marshal(feedback)
		if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "feedback",
			Operation: "CREATE",
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
