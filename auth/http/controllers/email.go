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
		email := c.Get("auth_user_email").(string)
		proposalStatus := c.QueryParam("status")

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
		m.SetHeader("To", email)

		if proposalStatus == "submitted" {
			m.SetHeader("Subject", "New Proposal on your Request")
			m.SetBody("text/html", "<p>Somebody just submitted a new proposal on your ride request. Check it out in the app.</p>")
		} else if proposalStatus == "accepted" {
			m.SetHeader("Subject", "Proposal accepted for your Request")
			m.SetBody("text/html", "<p>Somebody just accepted proposal on your ride request. Check it out in the app.</p>")
		} else if proposalStatus == "rejected" {
			m.SetHeader("Subject", "Proposal rejected for your Request")
			m.SetBody("text/html", "<p>Somebody just rejected proposal on your ride request. Check it out in the app.</p>")
		}

		if err := d.DialAndSend(m); err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		return cr.SendSuccessResponse(&c)
	}
}

