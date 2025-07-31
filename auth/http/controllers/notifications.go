package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

func SendNotification(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Extract the notification details from the request
		type Body struct {
			Title   string 	`json:"title"`
			To 		string 	`json:"to"`
			Body 	string 	`json:"body"`
		}

		var body Body
		if err := cr.BindAndValidate(&body, &c); err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}
		
		reqBody := map[string]any{
			"to":       body.To,
			"title":    body.Title,
			"body":     body.Body,
		}

		jsonData, err := json.Marshal(reqBody)
		if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		// send http request to external API
		req, err := http.NewRequest("POST", "https://exp.host/--/api/v2/push/send", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Accept-Encoding", "gzip, deflate")
		req.Header.Set("host", "exp.host")
	
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			cr.APIResponse.Error = nil
			return cr.SendSuccessResponse(&c)
		}
		defer resp.Body.Close()

		return cr.SendSuccessResponse(&c)
	}
}
