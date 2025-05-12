package controllers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func GetLocationName(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		lat := c.QueryParam("lat")
		long := c.QueryParam("long")

		// make a request to the Google Maps API
		apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
		url := "https://maps.googleapis.com/maps/api/geocode/json?latlng=" + lat + "," + long + "&key=" + apiKey

		resp, err := http.Get(url)
		if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			cr.APIResponse.Error = "Failed to get location name"
			return cr.SendErrorResponse(&c)
		}
		var result map[string]any
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}
		if result["status"] != "OK" {
			cr.APIResponse.Error = "Failed to get location name"
			return cr.SendErrorResponse(&c)
		}
		addressComponents := result["results"].([]any)[0].(map[string]any)
		address := addressComponents["formatted_address"].(string)
		if address == "" {
			cr.APIResponse.Error = "Failed to get location name"
			return cr.SendErrorResponse(&c)
		}
		cr.APIResponse.Data = map[string]any{
			"location_name": address,
		}
		cr.APIResponse.Status = "Successfully fetched location name"
		cr.APIResponse.Error = nil

		return cr.SendSuccessResponse(&c)
	}
}
