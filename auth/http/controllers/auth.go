package controllers

import (
	"fmt"
	"strings"

	"github.com/IbraheemHaseeb7/fyp-backend/auth"
	"github.com/IbraheemHaseeb7/fyp-backend/utils"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/labstack/echo/v4"
)

func Signup(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		type Request struct {
			Name               string `json:"name" validate:"required"`
			DOB                string `json:"dob" validate:"required"`
			Password           string `json:"password" validate:"required,password"`
			StudentCardURI     string `json:"studentCardURI"`
			LivePictureURI     string `json:"livePictureURI"`
			RegistrationNumber string `json:"registrationNumber" validate:"reg-no"`
			Department         string `json:"department"`
			Semester           int8   `json:"semester" validate:"required,gte=1,lte=12"`
			Email              string `json:"email" validate:"required,email"`
		}
		var reqBody Request
		if err := cr.BindAndValidate(&reqBody, &c); err != nil {
			return cr.SendErrorResponse(&c)
		}

		// hashing password
		hashedPwd, err := auth.HashPassword(reqBody.Password)
		if err != nil {
			return cr.SendErrorResponse(&c)
		}
		reqBody.Password = hashedPwd

		uuid := watermill.NewUUID()
		utils.Requests[uuid] = make(chan pubsub.PubsubMessage)

		// publishing a create message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "users",
			Operation: "CREATE",
			Topic:     "auth->db",
			UUID:      uuid,
			Payload: fmt.Sprintf(`
			{
				"name": "%s",
				"email": "%s",
				"registrationNumber": "%s",
				"dob": "%s",
				"semester": %d,
				"department": "%s",
				"password": "%s"
			}`, reqBody.Name, reqBody.Email, reqBody.RegistrationNumber,
				strings.Split(reqBody.DOB, "T")[0], reqBody.Semester, reqBody.Department,
				reqBody.Password),
		}
		err = cr.Publisher.PublishMessage(pubsubMessage)
		if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		cr.GetAndFormResponse(pubsubMessage)

		if cr.APIResponse.Error == "" || cr.APIResponse.Error == nil {

			v, ok := cr.APIResponse.Data.(map[string]any)
			if !ok {
				return cr.SendErrorResponse(&c)
			}

			id, ok := v["id"].(float64)
			if !ok {
				return cr.SendErrorResponse(&c)
			}

			acccessToken, err := auth.CreateToken(*auth.NewUserToken(id, reqBody.Name, reqBody.Email, reqBody.RegistrationNumber), 60)
			refreshToken, err := auth.CreateToken(*auth.NewUserToken(id, reqBody.Name, reqBody.Email, reqBody.RegistrationNumber), 1440)
			if err != nil {
				return cr.SendErrorResponse(&c)
			}
			response := map[string]string{
				"accessToken":  acccessToken,
				"refreshToken": refreshToken,
			}

			cr.APIResponse.Data = response
		}
		return cr.SendResponse(&c)
	}
}

func Login(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		type Request struct {
			Password           string `json:"password" validate:"required,password"`
			RegistrationNumber string `json:"registrationNumber" validate:"reg-no"`
		}
		var reqBody Request
		if err := cr.BindAndValidate(&reqBody, &c); err != nil {
			return cr.SendErrorResponse(&c)
		}

		uuid := watermill.NewUUID()
		utils.Requests[uuid] = make(chan pubsub.PubsubMessage)

		// publishing a create message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "users",
			Operation: "READ_ONE",
			Topic:     "auth->db",
			UUID:      uuid,
			Payload: fmt.Sprintf(`
			{
				"registrationNumber": "%s"
			}`, reqBody.RegistrationNumber),
		}
		if err := cr.Publisher.PublishMessage(pubsubMessage); err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		cr.GetAndFormResponse(pubsubMessage)

		// verifying password
		if !auth.CheckPasswordHash(reqBody.Password, cr.APIResponse.Data.(map[string]any)["password"].(string)) {
			cr.APIResponse.Status = "Invalid email or password"
			return cr.SendErrorResponse(&c)
		}

		if cr.APIResponse.Error == "" || cr.APIResponse.Error == nil {

			email := cr.APIResponse.Data.(map[string]any)["email"].(string)
			name := cr.APIResponse.Data.(map[string]any)["name"].(string)
			id := cr.APIResponse.Data.(map[string]any)["id"].(float64)

			acccessToken, err := auth.CreateToken(*auth.NewUserToken(id, name, email, reqBody.RegistrationNumber), 60)
			refreshToken, err := auth.CreateToken(*auth.NewUserToken(id, name, email, reqBody.RegistrationNumber), 1440)
			if err != nil {
				return cr.SendErrorResponse(&c)
			}
			response := map[string]string{
				"accessToken":  acccessToken,
				"refreshToken": refreshToken,
			}

			cr.APIResponse.StatusCode = 201
			cr.APIResponse.Data = response
		}
		return cr.SendResponse(&c)
	}
}

func Me(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		regNo := c.Get("auth_user_registration_number").(string)

		uuid := watermill.NewUUID()
		utils.Requests[uuid] = make(chan pubsub.PubsubMessage)

		// publishing a read message
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "users",
			Operation: "READ_ONE",
			Topic:     "auth->db",
			UUID:      uuid,
			Payload:   `{"registrationNumber": "` + regNo + `"}`,
		}
		err := cr.Publisher.PublishMessage(pubsubMessage)
		if err != nil {
			return cr.SendErrorResponse(&c)
		}

		cr.GetAndFormResponse(pubsubMessage)
		return cr.SendResponse(&c)
	}
}

func RefreshToken(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		type Request struct {
			RefreshToken string `json:"refreshToken" validate:"required"`
		}
		var reqBody Request
		if err := cr.BindAndValidate(&reqBody, &c); err != nil {
			return cr.SendErrorResponse(&c)
		}

		accessToken, refreshToken, err := auth.RefreshToken(reqBody.RefreshToken)
		if err != nil {
			return cr.SendErrorResponse(&c)
		}

		cr.APIResponse.Data = map[string]string{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		}
		cr.APIResponse.StatusCode = 200
		return cr.SendResponse(&c)
	}
}
