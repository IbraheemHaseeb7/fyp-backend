package controllers

import (
	"myapp/utils"

	"github.com/labstack/echo/v4"
)

func Login(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		// forming user request for the db service
		type RequestBody struct {
			Email string `json:"email" validate:"required,email"`
			Password string `json:"password" validate:"required,min=8"`
		}
		var reqBody RequestBody
		MapRequest(c, &reqBody)

		// Validate the struct using the validator instance
		if err := cr.Validate.Struct(reqBody); err != nil { return SendValidationError(c) }

		// generating JWT token
		token, err := utils.CreateToken(*utils.NewUserToken("IbraheemHaseeb7", reqBody.Email))
		utils.ErrorHandler(err)

		response := GenerateJsonResponse("{\"token\": \"" + token + "\"}", false)

		return SendCustomSuccessResponse(c, &response)		 
	}
}

func Signup(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		// forming user request for the db service
		type RequestBody struct {
			Username string `json:"username" validate:"required,min=3,max=30"`
			Email string `json:"email" validate:"required,email"`
			Password string `json:"password" validate:"required,min=8"`
		}
		var reqBody RequestBody
		MapRequest(c, &reqBody)

		// Validate the struct using the validator instance
		if err := cr.Validate.Struct(reqBody); err != nil { return SendValidationError(c) }

		// generating JWT token
		token, err := utils.CreateToken(*utils.NewUserToken("IbraheemHaseeb7", reqBody.Email))
		utils.ErrorHandler(err)

		response := GenerateJsonResponse("{\"token\": \"" + token + "\"}", false)

		return SendCustomSuccessResponse(c, &response)		 
	}
}
