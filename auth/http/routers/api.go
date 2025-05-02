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
	router.PATCH("/me", controllers.UpdateUser(controllers.ControllerRequest{
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
	}), middlewares.Authenticate())
	router.GET("/vehicle/:id", controllers.ReadOneVehicle(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())
	router.POST("/vehicle", controllers.CreateVehicle(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())
	router.PATCH("/vehicle/:id", controllers.UpdateVehicle(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())
	router.DELETE("/vehicle/:id", controllers.DeleteVehicle(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())

	/*
	 * EMAILS
	 *
	 * This contains all the routes for the emails
	 *
	 * - GET /email			sends email just for testing right now
	*/
	router.GET("/email", controllers.SendEmail(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())
	router.GET("/otp", controllers.SendOTP(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())
	router.POST("/verify-otp", controllers.VerifyOTP(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())

	/*
		* REQUESTS
		*
		* This contains all the routes for the requests
		*
		* - GET /requests			fetches all the requests from the database
		* - GET /requests/:id		fetches all the requests from the database
		* - POST /requests			fetches all the requests from the database
		* - PATCH /requests/:id		fetches all the requests from the database
		* - DELETE /requests/:id	fetches all the requests from the database
	*/
	router.GET("/requests", controllers.GetAllRequests(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())
	router.GET("/requests/:id", controllers.GetSingleRequest(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())
	router.POST("/requests", controllers.CreateRequest(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())
	router.PATCH("/requests/:id", controllers.UpdateRequest(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())
	router.DELETE("/requests/:id", controllers.DeleteRequest(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())
	router.POST("/requests/set-status", controllers.SetStatus(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())


	router.GET("/proposals/:request_id", controllers.GetAllProposals(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())
	router.GET("/proposals/user/me", controllers.GetAllMyProposals(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())
}
