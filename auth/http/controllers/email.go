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
		type Request struct {
			Email string `json:"email" validate:"required,cui-email"`
			Subject string `json:"subject" validate:"required"`
			Body   string `json:"body" validate:"required"`
		}
		var reqBody Request
		if err := cr.BindAndValidate(&reqBody, &c); err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

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
		m.SetHeader("To", reqBody.Email)

		m.SetHeader("Subject", reqBody.Subject)
		m.SetBody("text/html", reqBody.Body)

		if err := d.DialAndSend(m); err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		return cr.SendSuccessResponse(&c)
	}
}

