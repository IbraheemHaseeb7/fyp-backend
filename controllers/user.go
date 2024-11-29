package controllers

import (
	"myapp/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

func UsersCreate(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		// forming user request for the db service
		type RequestBody struct {
			Name string `json:"name" validate:"required,min=3,max=32"`
			Age uint8 `json:"age" validate:"required,gte=1,lte=100"`
		}
		var reqBody RequestBody
		MapRequest(c, &reqBody)

		// Validate the struct using the validator instance
		if err := cr.Validate.Struct(reqBody); err != nil { return SendValidationError(c) }

		// creating payload for db service
		payload := GeneratePayload(Payload{
			Entity: "users",
			Operation: "create",
			Payload: reqBody,
		})

		// sending and receiving request from db service
		result := GetAndSendRequest(&cr, payload)

		return SendCustomSuccessResponse(c, &result)
	}
}

func UsersRead(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		limit, offset := utils.GetLimitAndOffset(c.QueryParam("page"))
		type Request struct {
			Limit int
			Offset int
		}
		var reqBody Request
		reqBody.Limit = limit
		reqBody.Offset = offset
	
		// generating payload for the db service
		payload := GeneratePayload(Payload{
			Entity: "users",
			Operation: "read",
			Payload: reqBody,
		})

		// sending and receiving request from db service
		result := GetAndSendRequest(&cr, payload)
		response := GenerateJsonResponse(result, true)

		return SendCustomSuccessResponse(c, response)
	}
}

func UsersUpdate(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		// forming user request for the db service
		userId := utils.StrToInt(c.QueryParam("id"))
		type RequestBody struct {
			Name string `json:"name"`
			Age int `json:"age"`
			Id int 
		}
		var reqBody RequestBody
		MapRequest(c, &reqBody)

		// sending user id from params into the request
		reqBody.Id = userId

		// generating payload for the db service
		payload := GeneratePayload(Payload{
			Entity: "users",
			Operation: "update",
			Payload: reqBody,
		})

		result := GetAndSendRequest(&cr, payload) 
		response := GenerateJsonResponse(result, false)
		
		return SendCustomSuccessResponse(c, &response)
	}
}

func UsersDelete(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		type Request struct {
			Id int
		}
		// forming user request
		userId := c.QueryParam("id")
		var reqBody Request
		convertedUserId, err := strconv.Atoi(userId)
		utils.ErrorHandler(err)
		reqBody.Id = convertedUserId

		payload := GeneratePayload(Payload{
			Entity: "users",
			Operation: "delete",
			Payload: reqBody,	
		})

		// sending and receiving request from db service
		result := GetAndSendRequest(&cr, payload)

		return SendCustomSuccessResponse(c, &result)
	}
}

