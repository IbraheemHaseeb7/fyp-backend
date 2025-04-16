package controllers

import (
	"crypto/tls"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"gopkg.in/gomail.v2"
)

func SendEmail(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {

		SMTP_HOST := os.Getenv("SMTP_HOST")
		port := os.Getenv("SMTP_PORT")
		SMTP_EMAIL := os.Getenv("SMTP_EMAIL")
		SMTP_PASSWORD := os.Getenv("SMTP_PASSWORD")

		SMTP_PORT, err := strconv.Atoi(port)
		if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		d := gomail.NewDialer(SMTP_HOST, SMTP_PORT, "ibraheemhaseeb7", SMTP_PASSWORD)
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

		m := gomail.NewMessage()
		m.SetHeader("From", SMTP_EMAIL)
		m.SetHeader("To", "fa21-bcs-052@cuilahore.edu.pk")
		m.SetHeader("Subject", "Hello from Ridelink")
		m.SetBody("text/html", "<p>How are you my pookie lil cutie patootie???</p>")

		if err := d.DialAndSend(m); err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		return cr.SendSuccessResponse(&c)
	}
}
