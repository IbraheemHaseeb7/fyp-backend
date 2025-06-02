package controllers

import (
	"encoding/json"

	"github.com/IbraheemHaseeb7/fyp-backend/utils"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/labstack/echo/v4"
)

func GetAllProposals(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {

		userId := c.Get("auth_user_id")
		requestId := c.Param("request_id")

		payload, err := json.Marshal(map[string]any{
			"request_id": requestId,
			"user_id": userId,
		}); if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		uuid := watermill.NewUUID()
		utils.Requests.Store(uuid, make(chan pubsub.PubsubMessage))

		// publishing a read message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "proposals",
			Operation: "GET_ALL",
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

func GetAllMyProposals(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Get("auth_user_id")

		payload, err := json.Marshal(map[string]any{
			"user_id": userId,
		}); if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		uuid := watermill.NewUUID()
		utils.Requests.Store(uuid, make(chan pubsub.PubsubMessage))

		// publishing a read message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "proposals",
			Operation: "GET_ALL_MY",
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
