package controllers

import (
	"fmt"
	"strconv"

	"github.com/IbraheemHaseeb7/fyp-backend/utils"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/labstack/echo/v4"
)

func ReadOneVehicle(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		regNo := c.Get("auth_user_registration_number").(string)

		uuid := watermill.NewUUID()
		utils.Requests.Store(uuid, make(chan pubsub.PubsubMessage))

		// publishing a read message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "vehicles",
			Operation: "READ_ONE",
			Topic:     "auth->db",
			UUID:      uuid,
			Payload:   `{"id": ` + id + `, "regNo": "` + regNo + `"}`,
		}
		err := cr.Publisher.PublishMessage(pubsubMessage)
		if err != nil {
			return cr.SendErrorResponse(&c)
		}

		cr.GetAndFormResponse(pubsubMessage)

		v, ok := cr.APIResponse.Data.(map[string]any)
		if !ok {
			return cr.SendErrorResponse(&c)
		}

		if v, ok := v["id"].(float64); !ok || v == 0 {
			cr.APIResponse.StatusCode = 404
			return cr.SendErrorResponse(&c)
		}
		return cr.SendResponse(&c)
	}
}

func ReadAllVehicles(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {

		regNo, ok := c.Get("auth_user_registration_number").(string)
		if !ok {
			return cr.SendErrorResponse(&c)
		}

		uuid := watermill.NewUUID()
		utils.Requests.Store(uuid, make(chan pubsub.PubsubMessage))
		page := c.QueryParam("page")

		// publishing a read message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "vehicles",
			Operation: "READ_ALL",
			Topic:     "auth->db",
			UUID:      uuid,
			Payload:   `{"regNo": "` + regNo + `", "pageNo": ` + page + `}`,
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
		id, ok := c.Get("auth_user_id").(float64)
		if !ok {
			return cr.SendErrorResponse(&c)
		}

		type Request struct {
			Make      string  `json:"make" validate:"required"`
			Model     string  `json:"model" validate:"required"`
			Year      int     `json:"year" validate:"required"`
			Type      string  `json:"type" validate:"oneof=bike car rikshaw"`
			UserID    float64 `json:"userId"`
			VIN       string  `json:"vin" validate:"required"`
			FrontURI  string  `json:"frontUri" validate:"required"`
			BackURI   string  `json:"backUri" validate:"required"`
			InsideURI string  `json:"insideUri" validate:"required"`
		}
		var reqBody Request
		if err := cr.BindAndValidate(&reqBody, &c); err != nil {
			return cr.SendErrorResponse(&c)
		}

		uuid := watermill.NewUUID()
		utils.Requests.Store(uuid, make(chan pubsub.PubsubMessage))

		reqBody.UserID = id

		// publishing a read message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "vehicles",
			Operation: "CREATE",
			Topic:     "auth->db",
			UUID:      uuid,
			Payload: fmt.Sprintf(`{
				"make": "%s",
				"model": "%s",
				"year": %d,
				"type": "%s",
				"userId": %.f,
				"vin": "%s",
				"backUri": "%s",
				"frontUri": "%s",
				"insideUri": "%s"
				}`, reqBody.Make, reqBody.Model, reqBody.Year, reqBody.Type, reqBody.UserID,
				reqBody.VIN, reqBody.BackURI, reqBody.FrontURI, reqBody.InsideURI),
		}
		err := cr.Publisher.PublishMessage(pubsubMessage)
		if err != nil {
			return cr.SendErrorResponse(&c)
		}

		cr.GetAndFormResponse(pubsubMessage)
		return cr.SendResponse(&c)
	}
}

func UpdateVehicle(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, ok := c.Get("auth_user_id").(float64)
		if !ok {
			return cr.SendErrorResponse(&c)
		}

		type Request struct {
			Make      string  `json:"make" validate:""`
			Model     string  `json:"model" validate:"required"`
			Year      int     `json:"year" validate:"required"`
			Type      string  `json:"type" validate:"oneof=bike car rikshaw"`
			UserID    float64 `json:"userId"`
			ID        int     `json:"id"`
			VIN       string  `json:"vin" validate:"required"`
			FrontURI  string  `json:"frontUri" validate:"required"`
			BackURI   string  `json:"backUri" validate:"required"`
			InsideURI string  `json:"insideUri" validate:"required"`
		}
		var reqBody Request
		if err := cr.BindAndValidate(&reqBody, &c); err != nil {
			return cr.SendErrorResponse(&c)
		}

		uuid := watermill.NewUUID()
		utils.Requests.Store(uuid, make(chan pubsub.PubsubMessage))

		reqBody.UserID = id
		vId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return cr.SendErrorResponse(&c)
		}
		reqBody.ID = vId

		// publishing a read message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "vehicles",
			Operation: "UPDATE",
			Topic:     "auth->db",
			UUID:      uuid,
			Payload: fmt.Sprintf(`{
				"make": "%s",
				"model": "%s",
				"year": %d,
				"type": "%s",
				"userId": %.f,
				"vin": "%s",
				"backUri": "%s",
				"frontUri": "%s",
				"insideUri": "%s",
				"id": %d
				}`, reqBody.Make, reqBody.Model, reqBody.Year, reqBody.Type, reqBody.UserID,
				reqBody.VIN, reqBody.BackURI, reqBody.FrontURI, reqBody.InsideURI, reqBody.ID),
		}
		err = cr.Publisher.PublishMessage(pubsubMessage)
		if err != nil {
			return cr.SendErrorResponse(&c)
		}

		cr.GetAndFormResponse(pubsubMessage)
		return cr.SendResponse(&c)
	}
}

func DeleteVehicle(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, ok := c.Get("auth_user_id").(float64)
		if !ok {
			return cr.SendErrorResponse(&c)
		}

		vId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return cr.SendErrorResponse(&c)
		}

		uuid := watermill.NewUUID()
		utils.Requests.Store(uuid, make(chan pubsub.PubsubMessage))

		// publishing a read message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "vehicles",
			Operation: "DELETE",
			Topic:     "auth->db",
			UUID:      uuid,
			Payload: fmt.Sprintf(`{
				"vId": %d,
				"id": %.f
				}`, vId, id),
		}
		err = cr.Publisher.PublishMessage(pubsubMessage)
		if err != nil {
			return cr.SendErrorResponse(&c)
		}

		cr.GetAndFormResponse(pubsubMessage)
		return cr.SendResponse(&c)
	}
}
