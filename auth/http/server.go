package http

import (
	"github.com/IbraheemHaseeb7/fyp-backend/http/routers"
	"github.com/IbraheemHaseeb7/fyp-backend/utils"
	"github.com/IbraheemHaseeb7/fyp-backend/validation"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func StartHTTPServer(p *pubsub.Publisher) {
	utils.Requests = make(map[string]chan pubsub.PubsubMessage)

	validation.Validate = validator.New(validator.WithRequiredStructEnabled())
	validation.Validate.RegisterValidation("reg-no", validation.IsComsatsRegistrationNumber)
	validation.Validate.RegisterValidation("password", validation.IsPassword)

	e := echo.New()

	routers.ApiRouter(e, p)

	e.Logger.Fatal(e.Start(":8000"))
}
