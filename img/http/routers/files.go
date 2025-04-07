package routers

import (
	"github.com/IbraheemHaseeb7/fyp-backend/http/controllers"
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
	* - GET /files			allows users to view files
	* - POST /files			allows the creation of new files that are compressed with a unique ID
	*/
	router.GET("/public/:level1/:level2/:filename", filesHandler.ReadOnePublicFile)
	router.GET("/private/:level1/:level2/:filename", controllers.ReadOnePrivateFile(controllers.ControllerRequest{
		Publisher: p,
	}))
	router.POST("/", controllers.UploadOneFile(controllers.ControllerRequest{
		Publisher: p,
	}))
	router.DELETE("/", controllers.DeleteFiles(controllers.ControllerRequest{
		Publisher: p,
	}))
}
