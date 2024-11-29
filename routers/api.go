package routers

import (
	"myapp/controllers"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/streadway/amqp"
)

func ApiRouter(e *echo.Echo, ch *amqp.Channel, reqQ *amqp.Queue, resQ *amqp.Queue, receiverChanMap map[string]chan string, validate *validator.Validate) {
	router := e.Group("/api")

	router.POST("/", controllers.UsersCreate(controllers.ControllerRequest{
		Ch: ch, ResQ: resQ, ReqQ: reqQ, ReceiverChanMap: receiverChanMap, Validate: validate,
	}))
	router.GET("/", controllers.UsersRead(controllers.ControllerRequest{
		Ch: ch, ResQ: resQ, ReqQ: reqQ, ReceiverChanMap: receiverChanMap, Validate: validate,
	}))
	router.PATCH("/", controllers.UsersUpdate(controllers.ControllerRequest{
		Ch: ch, ResQ: resQ, ReqQ: reqQ, ReceiverChanMap: receiverChanMap, Validate: validate,
	}))
	router.DELETE("/", controllers.UsersDelete(controllers.ControllerRequest{
		Ch: ch, ResQ: resQ, ReqQ: reqQ, ReceiverChanMap: receiverChanMap, Validate: validate,
	}))


	router.POST("/login", controllers.Login(controllers.ControllerRequest{
		Ch: ch, ResQ: resQ, ReqQ: reqQ, ReceiverChanMap: receiverChanMap, Validate: validate,
	}))
	router.POST("/signup", controllers.Signup(controllers.ControllerRequest{
		Ch: ch, ResQ: resQ, ReqQ: reqQ, ReceiverChanMap: receiverChanMap, Validate: validate,
	}))
}
