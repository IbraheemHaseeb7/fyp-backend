package controllers

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/IbraheemHaseeb7/fyp-backend/utils"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/IbraheemHaseeb7/types"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/labstack/echo/v4"
)

func GetAllRides(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Get("auth_user_id").(float64)

		uuid := watermill.NewUUID()
		utils.Requests[uuid] = make(chan pubsub.PubsubMessage)
		page, err := strconv.Atoi(c.QueryParam("page")); if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		payload, err := json.Marshal(map[string]any{"page": page, "user_id": userId}); if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		// publishing a read message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "rides",
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

func GetSingleRide(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Get("auth_user_id").(float64)

		uuid := watermill.NewUUID()
		utils.Requests[uuid] = make(chan pubsub.PubsubMessage)
		id := c.Param("id")

		payload, err := json.Marshal(map[string]any{"id": id, "user_id": userId}); if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		// publishing a read message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "rides",
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

func CreateRide(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		uuid := watermill.NewUUID()
		utils.Requests[uuid] = make(chan pubsub.PubsubMessage)

		type Ride struct {
			VehicleID int64 `json:"vehicle_id" validate:"required"`
			DriverID  int64 `json:"driver_id" validate:"required"`
			PassengerID int64 `json:"passenger_id" validate:"required"`
			RequestID int64 `json:"request_id" validate:"required"`
			ProposalID int64 `json:"proposal_id" validate:"required"`
			StartTime time.Time 
			EndTime   time.Time 
		}

		var reqBody Ride
		if err := cr.BindAndValidate(&reqBody, &c); err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		reqBody.StartTime = time.Now()
		reqBody.EndTime = time.Now().Add(time.Hour * 1)

		payload, err := json.Marshal(reqBody); if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		// publishing a read message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "rides",
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

func UpdateRide(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		uuid := watermill.NewUUID()
		utils.Requests[uuid] = make(chan pubsub.PubsubMessage)
		id := c.Param("id")

		var reqBody types.Ride
		if err := cr.BindAndValidate(&reqBody, &c); err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		payload, err := json.Marshal(map[string]any{"id": id, "data": reqBody}); if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		// publishing a read message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "rides",
			Operation: "UPDATE",
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

func DeleteRide(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Get("auth_user_id").(float64)

		uuid := watermill.NewUUID()
		utils.Requests[uuid] = make(chan pubsub.PubsubMessage)
		id := c.Param("id")

		payload, err := json.Marshal(map[string]any{"id": id, "user_id": userId}); if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		// publishing a read message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "rides",
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

// fetch user's active ride
func ActiveRide(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Get("auth_user_id")

		uuid := watermill.NewUUID()
		utils.Requests[uuid] = make(chan pubsub.PubsubMessage)

		payload, err := json.Marshal(map[string]any{"user_id": userId}); if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		// publishing a read message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "rides",
			Operation: "ACTIVE_RIDE",
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
