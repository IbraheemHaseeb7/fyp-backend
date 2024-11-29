package main

import (
	"myapp/routers"
	"myapp/utils"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	// fetching env variables
	godotenv.Load()
	
	// validator setup
	validate := validator.New(validator.WithRequiredStructEnabled())

	// starting rabbit mq
	receiverChanMap := make(map[string]chan string)
	ch, reqQ, resQ := utils.RabbitMQ(&receiverChanMap)

	// setting up echo server
	e := echo.New()

	// registering routers
	routers.ApiRouter(e, ch, &reqQ, &resQ, receiverChanMap, validate)

	e.Logger.Fatal(e.Start(":1323"))
}
