package routers

import (
	"github.com/IbraheemHaseeb7/fyp-backend/http/controllers"
	"github.com/IbraheemHaseeb7/fyp-backend/http/middlewares"
	"github.com/IbraheemHaseeb7/fyp-backend/http/services"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/labstack/echo/v4"
)

func ImgRouter(e *echo.Echo, p *pubsub.Publisher) {

	router := e.Group("/files")

	controllerRequest := controllers.ControllerRequest{
		Publisher: p,
	}
	filesServices := services.FilesServices{}
	filesHandler := controllers.FileHandler{
		F: &filesServices,
		Cr: &controllerRequest,
	}

	/*
	* FILES
	*
	* This contains all the routes for the files
	*
	* - GET /public/...			allows users to view files
	* - GET /private/...		allows users to view files
	* - POST /					allows the creation of new files that are compressed with a unique ID
	* - DELETE/					allows the deletion of files that are compressed with a unique ID
	* - DELETE/except			allows the deletion of files that are compressed with a unique ID
	*/
	router.GET("/public/:level1/:level2/:filename", filesHandler.ReadOnePublicFile)
	router.GET("/private/:level1/:level2/:filename", controllers.ReadOnePrivateFile(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate(controllers.ControllerRequest{Publisher: p}))
	router.POST("/", controllers.UploadOneFile(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate(controllers.ControllerRequest{Publisher: p}))
	router.DELETE("/", controllers.DeleteFiles(controllers.ControllerRequest{
		Publisher: p,
	}))
	router.DELETE("/except/", controllers.KeepOnlyFiles(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate(controllers.ControllerRequest{Publisher: p}))

	/*
	* VERIFY
	*
	* This contains all the routes for the verification of the files
	*
	* - POST /verify-card			allows the verification of the cards (particularly CUI cards only)
	*/
	router.POST("/verify-card", controllers.VerifyCard(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate(controllers.ControllerRequest{Publisher: p}))
	router.POST("/verify-selfie", controllers.VerifySelfie(controllers.ControllerRequest{
		Publisher: p,
	}), middlewares.Authenticate(controllers.ControllerRequest{Publisher: p}))
}
