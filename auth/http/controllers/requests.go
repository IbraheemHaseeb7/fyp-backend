package controllers

import (
	"encoding/json"
	"fmt"
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
		utils.Requests.Store(uuid, make(chan pubsub.PubsubMessage))
		page, err := strconv.Atoi(c.QueryParam("page")); if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		status := c.QueryParam("status"); if status == "" {
			status = "searching"
		}

		me := c.QueryParam("me"); if me == "true" {
			me = fmt.Sprintf("%.f", c.Get("auth_user_id").(float64))
		} else {
			me = "false"
		}

		payload, err := json.Marshal(map[string]any{"page": page, "status": status, "me": me}); if err != nil {
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
		utils.Requests.Store(uuid, make(chan pubsub.PubsubMessage))
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
		utils.Requests.Store(uuid, make(chan pubsub.PubsubMessage))

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
		type Request struct {
			Trunk			*bool			`json:"trunk"`
			VehicleID		int64			`json:"vehicle_id"`
			RequestID		int64			`json:"request_id"`
			Persons 		uint8			`json:"persons" validate:""`
			VehicleType		string			`json:"vehicle_type" validate:""`
			FromLat			float64			`json:"from_lat" validate:""`
			FromLong		float64			`json:"from_long" validate:""`
			ToLat			float64			`json:"to_lat" validate:""`
			ToLong			float64			`json:"to_long" validate:""`
			Status			string			`json:"status" validate:""`
			OriginatorRole	string			`json:"originator_role" validate:""`
		}
		var reqBody Request
		if err := cr.BindAndValidate(&reqBody, &c); err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		mapData, err := utils.StructToMap(reqBody); if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		for key, value := range mapData {
			if key == "trunk" {
				if (value.(*bool)) == nil {
					delete(mapData, "trunk")
				}
			}
			if value == 0 || value == "0" || value == "" || value == nil || value == float64(0) || value == uint8(0) || value == int64(0)  {
				delete(mapData, key)
			}
		}

		mapData["user_id"] = c.Get("auth_user_id")
		mapData["id"] = c.Param("id")
		payload, err := json.Marshal(mapData); if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		uuid := watermill.NewUUID()
		utils.Requests.Store(uuid, make(chan pubsub.PubsubMessage))

		// publishing a read message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    	"requests",
			Operation: 	"UPDATE",
			Topic:     	"auth->db",
			UUID:     	uuid,
			Payload: 	string(payload),
		}
		err = cr.Publisher.PublishMessage(pubsubMessage)
		if err != nil {
			return cr.SendErrorResponse(&c)
		}

		cr.GetAndFormResponse(pubsubMessage)
		return cr.SendResponse(&c)
	}
}

func DeleteRequest(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {

		uuid := watermill.NewUUID()
		utils.Requests.Store(uuid, make(chan pubsub.PubsubMessage))
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

func SetStatus(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		type Request struct {
			ID 	int	`json:"id" validate:"required"`
			Status 	string	`json:"status" validate:"required"`
		}
		var reqBody Request
		if err := cr.BindAndValidate(&reqBody, &c); err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		mapData, err := utils.StructToMap(reqBody); if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}
		payload, err := json.Marshal(mapData); if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		uuid := watermill.NewUUID()
		utils.Requests.Store(uuid, make(chan pubsub.PubsubMessage))

		// publishing a read message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "requests",
			Operation: "SET_STATUS",
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

func GetMatchedProposalOfARequest(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {

		uuid := watermill.NewUUID()
		utils.Requests.Store(uuid, make(chan pubsub.PubsubMessage))
		id, err := strconv.Atoi(c.Param("id")); if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		payload, err := json.Marshal(map[string]any{"request_id": id}); if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		// publishing a read message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "requests",
			Operation: "GET_MATCHED_PROPOSAL_OF_A_REQUEST",
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

func GetMyProposalForARequest(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Get("auth_user_id").(float64)

		uuid := watermill.NewUUID()
		utils.Requests.Store(uuid, make(chan pubsub.PubsubMessage))
		id, err := strconv.Atoi(c.Param("id")); if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		payload, err := json.Marshal(map[string]any{"request_id": id, "user_id": userID}); if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		// publishing a read message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "requests",
			Operation: "GET_MY_PROPOSAL_FOR_A_REQUEST",
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

func GetMatches(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {

		type Request struct {
			FromLat			float64	`json:"from_lat" validate:"required"`
			FromLong 		float64	`json:"from_long" validate:"required"`
			ToLat			float64	`json:"to_lat" validate:"required"`
			ToLong			float64	`json:"to_long" validate:"required"`
			Id				float64	`json:"id" validate:"required"`
			OriginatorRole 	string	`json:"originator_role" validate:"required"`
			VehicleType 	string	`json:"vehicle_type" validate:"required"`
		}

		var reqBody Request
		if err := cr.BindAndValidate(&reqBody, &c); err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		uuid := watermill.NewUUID()
		utils.Requests.Store(uuid, make(chan pubsub.PubsubMessage))

		payload, err := json.Marshal(map[string]any{
			"id": reqBody.Id,
			"from_lat":  reqBody.FromLat,
			"from_long": reqBody.FromLong,
			"to_lat":    reqBody.ToLat,
			"to_long":   reqBody.ToLong,
			"originator_role": reqBody.OriginatorRole,
			"vehicle_type": reqBody.VehicleType,
		}); if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		// publishing a read message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "requests",
			Operation: "GET_MATCHES",
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

func GetActiveRequest(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {

		userId := c.Get("auth_user_id").(float64)

		uuid := watermill.NewUUID()
		utils.Requests.Store(uuid, make(chan pubsub.PubsubMessage))

		payload, err := json.Marshal(map[string]any{
			"user_id": userId,
		}); if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		// publishing a read message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "requests",
			Operation: "GET_ACTIVE_REQUEST",
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
