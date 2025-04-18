package controllers

import "github.com/labstack/echo/v4"

func GetAllRequests(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		
		return cr.SendResponse(&c)
	}
}

func GetSingleRequest(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		
		return cr.SendResponse(&c)
	}
}

func CreateRequest(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		
		return cr.SendSuccessResponse(&c)
	}
}

func UpdateRequest(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		
		return cr.SendSuccessResponse(&c)
	}
}

func DeleteRequest(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		
		return cr.SendSuccessResponse(&c)
	}
}
