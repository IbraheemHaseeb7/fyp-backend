package http

import (
	"github.com/IbraheemHaseeb7/fyp-backend/http/config"
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

	e := echo.New()

	config.ModifyDefaultResponses(e)

	routers.ImgRouter(e, p)

	e.Logger.Fatal(e.Start("192.168.1.8:8001"))
}
