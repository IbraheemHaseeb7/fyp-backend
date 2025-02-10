package routers

import (
	"github.com/IbraheemHaseeb7/fyp-backend/http/controllers"
	"github.com/IbraheemHaseeb7/fyp-backend/http/middlewares"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/labstack/echo/v4"
)

func ApiRouter(e *echo.Echo, p *pubsub.Publisher) {

	router := e.Group("/api")
	router.Use(middlewares.Headers())

	/*
			 * AUTH
			 *
		     * This contains all the routes for the auth
			 *
			 * - POST /signup	allows the creation of new users and generates JWT tokens
	*/
	router.POST("/signup", controllers.Signup(controllers.ControllerRequest{
		Publisher: p,
	}))
	router.POST("/login", controllers.Login(controllers.ControllerRequest{
		Publisher: p,
	}))
	router.POST("/refresh", controllers.RefreshToken(controllers.ControllerRequest{
		Publisher: p,
	}))
	router.GET("/me", controllers.Me(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())

	/*
			 * USERS
			 *
		     * This contains all the routes for the users
			 *
			 * - GET /user 		gives all the users with the ability to perform pagination
			 * - GET /user/:id	gives information about one user in particular
			 * - POST /user/	allows the creation of new users
	*/
	router.GET("/user", controllers.ReadAllUsers(controllers.ControllerRequest{
		Publisher: p,
	}))
	router.GET("/user/:id", controllers.ReadOneUser(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())
	router.POST("/user", controllers.CreateUsers(controllers.ControllerRequest{
		Publisher: p,
	}))

	/*
			 * VEHICLES
			 *
		     * This contains all the routes for the vehicles
			 *
			 * - GET /vechile 		gives all the vehicles with the ability to perform pagination
			 * - GET /vehicle/:id	gives information about one vehicle in particular
			 * - POST /vehicle/		allows the creation of new vehicles that could be assigned to a user
	*/
	router.GET("/vehicle", controllers.ReadAllVehicles(controllers.ControllerRequest{
		Publisher: p,
	}))
	router.GET("/vehicle/:id", controllers.ReadOneVehicle(controllers.ControllerRequest{
		Publisher: p,
	}))
	router.POST("/vehicle", controllers.CreateVehicle(controllers.ControllerRequest{
		Publisher: p,
	}))
}
