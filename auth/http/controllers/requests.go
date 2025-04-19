package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/IbraheemHaseeb7/fyp-backend/utils"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/IbraheemHaseeb7/types"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/labstack/echo/v4"
)

func GetAllRequests(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		uuid := watermill.NewUUID()
		utils.Requests[uuid] = make(chan pubsub.PubsubMessage)
		page, err := strconv.Atoi(c.QueryParam("page")); if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		payload, err := json.Marshal(map[string]int{"page": page}); if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		// publishing a read message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "requests",
			Operation: "READ_ALL",
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

func GetSingleRequest(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {

		uuid := watermill.NewUUID()
		utils.Requests[uuid] = make(chan pubsub.PubsubMessage)
		id, err := strconv.Atoi(c.Param("id")); if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		payload, err := json.Marshal(map[string]int{"id": id}); if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		// publishing a read message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "requests",
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

func CreateRequest(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {

		var reqBody types.Request

		reqBody.UserID = int64(c.Get("auth_user_id").(float64))

		if err := cr.BindAndValidate(&reqBody, &c); err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		uuid := watermill.NewUUID()
		utils.Requests[uuid] = make(chan pubsub.PubsubMessage)

		payload, err := json.Marshal(reqBody); if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		// publishing a read message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "requests",
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

func UpdateRequest(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		
		return cr.SendSuccessResponse(&c)
	}
}

func DeleteRequest(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {

		uuid := watermill.NewUUID()
		utils.Requests[uuid] = make(chan pubsub.PubsubMessage)
		id, err := strconv.Atoi(c.Param("id")); if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		payload, err := json.Marshal(map[string]int{"id": id, "user_id": int(c.Get("auth_user_id").(float64))}); if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		// publishing a read message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "requests",
			Operation: "DELETE",
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
