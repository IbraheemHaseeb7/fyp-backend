package routers

import (
	"github.com/IbraheemHaseeb7/fyp-backend/http/sockets"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/labstack/echo/v4"
)


func IndexRouter(e *echo.Echo, p *pubsub.Publisher) {

	router := e.Group("")

	/*
		* SOCKETS
		*
		* This contains all the routes for the sockets
		*
		* - GET /socket			starts the socket server
	*/
	socketServer := sockets.SetupSocket(p)
	router.GET("/socket.io/", echo.WrapHandler(socketServer))
	router.POST("/socket.io/", echo.WrapHandler(socketServer))
	router.GET("/", func(c echo.Context) error {
		return c.String(200, "Welcome to the socket server")
	})
}
