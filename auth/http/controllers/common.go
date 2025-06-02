package controllers

import (
	"github.com/IbraheemHaseeb7/fyp-backend/utils"
	"github.com/IbraheemHaseeb7/fyp-backend/validation"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/labstack/echo/v4"
)

type ControllerRequest struct {
	Publisher   *pubsub.Publisher
	APIResponse APIResponse
}

type APIResponse struct {
	Status     string `json:"status"`
	Error      any    `json:"error"`
	Data       any    `json:"data"`
	StatusCode int
}

func (cr *ControllerRequest) SendSuccessResponse(c *echo.Context) error {
	if cr.APIResponse.StatusCode == 0 {
		cr.APIResponse.StatusCode = 200
	}

	if cr.APIResponse.Status == "" {
		cr.APIResponse.Status = "Successfully processed request"
	}

	return (*c).JSON(cr.APIResponse.StatusCode, map[string]any{
		"status": cr.APIResponse.Status,
		"error":  cr.APIResponse.Error,
		"data":   cr.APIResponse.Data,
	})
}

func (cr *ControllerRequest) SendErrorResponse(c *echo.Context) error {
	if cr.APIResponse.StatusCode == 0 || (cr.APIResponse.StatusCode >= 200 && cr.APIResponse.StatusCode < 300) {
		cr.APIResponse.StatusCode = 500
	}

	if cr.APIResponse.Status == "" {
		cr.APIResponse.Status = "Error processing request"
	}

	cr.APIResponse.Data = nil

	return (*c).JSON(cr.APIResponse.StatusCode, map[string]any{
		"status": cr.APIResponse.Status,
		"error":  cr.APIResponse.Error,
		"data":   cr.APIResponse.Data,
	})
}

func (cr *ControllerRequest) SendResponse(c *echo.Context) error {

	if cr.APIResponse.StatusCode == 0 {
		cr.APIResponse.StatusCode = 200
	}

	if cr.APIResponse.Error != nil {
		cr.APIResponse.StatusCode = 500
	}

	if cr.APIResponse.Status == "" {
		if cr.APIResponse.StatusCode >= 200 && cr.APIResponse.StatusCode <= 299 {
			cr.APIResponse.Status = "Successfully processed request"
		} else if cr.APIResponse.StatusCode >= 300 && cr.APIResponse.StatusCode <= 599 {
			cr.APIResponse.Status = "Error processing request"
		} else {
			cr.APIResponse.Status = "Unknown status"
		}
	}

	return (*c).JSON(cr.APIResponse.StatusCode, map[string]any{
		"status": cr.APIResponse.Status,
		"error":  cr.APIResponse.Error,
		"data":   cr.APIResponse.Data,
	})
}

func (cr *ControllerRequest) GetAndFormResponse(pubsubMessage pubsub.PubsubMessage) {
	response := (<-utils.Requests.Load(pubsubMessage.UUID)).Payload.(map[string]any)
	utils.Requests.Delete(pubsubMessage.UUID)

	if response["status"] == nil {
		response["status"] = ""
	}

	if response["error"] == nil {
		if pubsubMessage.Operation == "CREATE" {
			cr.APIResponse.StatusCode = 201
		} else {
			cr.APIResponse.StatusCode = 200
		}
	} else {
		cr.APIResponse.StatusCode = 500
	}

	cr.APIResponse.Data = response["data"]
	cr.APIResponse.Error = response["error"]
	cr.APIResponse.Status = response["status"].(string)
}

func (cr *ControllerRequest) BindAndValidate(reqBody any, c *echo.Context) error {
	if err := (*c).Bind(reqBody); err != nil {
		cr.APIResponse.StatusCode = 500
		cr.APIResponse.Status = "Error processing request"
		cr.APIResponse.Error = err.Error()
		return err
	}

	if err := validation.Validate.Struct(reqBody); err != nil {
		cr.APIResponse.Error = err.Error()
		cr.APIResponse.Status = "Error processing request"
		cr.APIResponse.StatusCode = 400
		return err
	}

	return nil
}
