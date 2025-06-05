package http

import (
	"fmt"
	"os"

	"github.com/IbraheemHaseeb7/fyp-backend/http/config"
	"github.com/IbraheemHaseeb7/fyp-backend/http/routers"
	"github.com/IbraheemHaseeb7/fyp-backend/validation"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func StartHTTPServer(p *pubsub.Publisher) {
	validation.Validate = validator.New(validator.WithRequiredStructEnabled())

	e := echo.New()

	config.ModifyDefaultResponses(e)

	routers.ImgRouter(e, p)

	if os.Getenv("ENVIRONMENT") == "local" {
		e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
	} else if os.Getenv("ENVIRONMENT") == "staging" {
		e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", os.Getenv("BASE_ADDRESS"), os.Getenv("PORT"))))
	}
}
