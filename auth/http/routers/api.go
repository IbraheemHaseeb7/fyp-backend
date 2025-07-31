package routers

import (
	"github.com/IbraheemHaseeb7/fyp-backend/http/controllers"
	"github.com/IbraheemHaseeb7/fyp-backend/http/jobs"
	"github.com/IbraheemHaseeb7/fyp-backend/http/middlewares"
	"github.com/IbraheemHaseeb7/fyp-backend/utils"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/labstack/echo/v4"
)

func ApiRouter(e *echo.Echo, p *pubsub.Publisher) {

	router := e.Group("/api")
	router.Use(middlewares.Headers())

	poolSize := 5
	workerPool := jobs.NewWorkerPool(poolSize)
	workerPool.Start()

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
	router.POST("/reset-password", controllers.RefreshToken(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())
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
	router.POST("/store-device-token", controllers.StoreDeviceToken(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())

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
	 * - POST /email			sends email just for testing right now
	 * - GET /otp				sends otp to the user
	 * - POST /verify-otp		verifies the otp sent to the user
	*/
	router.POST("/email", controllers.SendEmail(controllers.ControllerRequest{
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
	router.GET("/matched-proposal/:id", controllers.GetMatchedProposalOfARequest(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())
	router.GET("/my-proposal/:id", controllers.GetMyProposalForARequest(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())
	router.POST("/find-matches", controllers.GetMatches(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())
	router.GET("/active-request", controllers.GetActiveRequest(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())

	/*
		* ROOMS
		*
		* This contains all the routes for the rooms
		*
		* - GET /rooms/:id		fetches all the rooms from the database
		*
	*/
	router.GET("/rooms/:id", controllers.GetRoom(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())
	router.GET("/rooms", controllers.GetRoomWithIDs(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())
	router.GET("/chats/:room_id", controllers.ChatMessages(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())

	/*
		* RIDES
		*
		* This contains all the routes for the rides
		*
		* - GET /rides			fetches all the rides from the database
		* - GET /rides/:id		fetches all the rides from the database
		* - POST /rides			fetches all the rides from the database
		* - PATCH /rides/:id		fetches all the rides from the database
		* - DELETE /rides/:id	fetches all the rides from the database
	*/
	router.GET("/rides", controllers.GetAllRides(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())
	router.GET("/rides/:id", controllers.GetSingleRide(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())
	router.POST("/rides", controllers.CreateRide(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())
	router.PATCH("/rides/:id", controllers.UpdateRide(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())
	router.DELETE("/rides/:id", controllers.DeleteRide(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())
	router.GET("/active-ride", controllers.ActiveRide(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())

	/*
		* PROPOSALS
		*
		* This contains all the routes for the proposals
		*
		* - GET /proposals/:request_id			fetches all the proposals from the database
		* - GET /proposals/user/me				fetches all the proposals from the database
	*/
	router.GET("/proposals/:request_id", controllers.GetAllProposals(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())
	router.GET("/proposals/user/me", controllers.GetAllMyProposals(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())


	/*
		* LOCATIONS
		*
		* This contains all the routes for the locations
		*
		* - GET /location			fetches the location from the google api
	*/
	router.GET("/location", controllers.GetLocationName(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())
	router.GET("/", func(c echo.Context) error {
		
		interalRequest := utils.NewInternalApiRequest("/api/email", "POST", 
			map[string]any{
				"email": "ibraheemibnhaseeb@gmail.com",
				"subject": "submitted",
				"body": "submitted",
			},
			map[string]string{
				"Content-Type": "application/json",
				"Authorization": c.Request().Header.Get("Authorization"),
			})
		task := jobs.NewTask(*interalRequest)
		workerPool.EnqueueTask(task)

		return c.String(200, "Welcome to the API")
	}, middlewares.Authenticate())

	/*
		* FEEDBACK
		*
		* This contains all the routes for the feedback
		*
		* - POST /feedback			fetches the feedback from the user
	*/
	router.POST("/feedback", controllers.AddFeedback(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())

	/*
		* NOTIFICATIONS
		*
		* This contains all the routes for the notifications
		*
		* - POST /notifications			sends push notifications to the user
	*/
	router.POST("/push-notification", controllers.SendNotification(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate())
}
