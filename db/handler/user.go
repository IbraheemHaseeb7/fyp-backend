package handler

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/IbraheemHaseeb7/fyp-backend/db"
	"github.com/IbraheemHaseeb7/fyp-backend/utils"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/IbraheemHaseeb7/types"
)

func Login(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {

	type Query struct {
		RegistrationNumber string `json:"registrationNumber"`
	}
	var query Query
	err := json.Unmarshal([]byte(pm.Payload.(string)), &query)
	if err != nil {
		return pubsub.PubsubMessage{
			Payload: map[string]any{
				"status": "Could not break down query",
				"error":  err.Error(),
			},
			Entity:    pm.Entity,
			Operation: pm.Operation,
			Topic:     "db->auth",
			UUID:      pm.UUID,
		}, nil
	}

	if query.RegistrationNumber == "" {
		return pubsub.PubsubMessage{
			Payload: map[string]any{
				"status": "Could not find registration number",
				"error":  "Please make sure to send registration number in the query",
			},
			Entity:    pm.Entity,
			Operation: pm.Operation,
			Topic:     "db->auth",
			UUID:      pm.UUID,
		}, nil
	}

	var user types.User
	db.DB.Model(&types.User{}).Where("registration_number = ?", query.RegistrationNumber).Find(&user)

	return pubsub.PubsubMessage{
		Payload: map[string]any{
			"data":   user,
			"status": "success",
			"error":  nil,
		},
		Entity:    pm.Entity,
		Operation: pm.Operation,
		Topic:     "db->auth",
		UUID:      pm.UUID,
	}, nil
}

func ReadOneUser(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {

	type Query struct {
		RegistrationNumber string `json:"registrationNumber"`
	}
	var query Query
	err := json.Unmarshal([]byte(pm.Payload.(string)), &query)
	if err != nil {
		return pubsub.PubsubMessage{
			Payload: map[string]any{
				"status": "Could not break down query",
				"error":  err.Error(),
			},
			Entity:    pm.Entity,
			Operation: pm.Operation,
			Topic:     "db->auth",
			UUID:      pm.UUID,
		}, nil
	}

	if query.RegistrationNumber == "" {
		return pubsub.PubsubMessage{
			Payload: map[string]any{
				"status": "Could not find registration number",
				"error":  "Please make sure to send registration number in the query",
			},
			Entity:    pm.Entity,
			Operation: pm.Operation,
			Topic:     "db->auth",
			UUID:      pm.UUID,
		}, nil
	}

	type RequestUser struct {
		Name 				string			`json:"name"`
		Email				string			`json:"email"`
		RegistrationNumber 	string			`json:"registrationNumber"`
		LivePictureURI		string			`json:"livePictureURI"`
		StudentCardURI		string			`json:"studentCardURI"`
		EmailVerifiedAt		sql.NullTime	`json:"emailVerifiedAt"`
		CardVerifiedAt		sql.NullTime	`json:"cardVerifiedAt"`
		SelfieVerifiedAt	sql.NullTime	`json:"selfieVerifiedAt"`
		Semester			uint8			`json:"semester"`
		Department			string			`json:"department"`
	}
	var user RequestUser
	db.DB.Model(&types.User{}).Where("registration_number = ?", query.RegistrationNumber).Find(&user)

	return pubsub.PubsubMessage{
		Payload: map[string]any{
			"data":   user,
			"status": "success",
			"error":  nil,
		},
		Entity:    pm.Entity,
		Operation: pm.Operation,
		Topic:     "db->auth",
		UUID:      pm.UUID,
	}, nil
}

func ReadAllUsers(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {

	type RequestUser struct {
		Name 				string			`json:"name"`
		Email				string			`json:"email"`
		RegistrationNumber 	string			`json:"registrationNumber"`
		LivePictureURI		string			`json:"livePictureURI"`
		StudentCardURI		string			`json:"studentCardURI"`
		EmailVerifiedAt		sql.NullTime	`json:"emailVerified"`
		CardVerifiedAt		sql.NullTime	`json:"cardVerified"`
		SelfieVerifiedAt	sql.NullTime	`json:"selfieVerified"`
		Semester			uint8			`json:"semester"`
		Department			string			`json:"department"`
	}
	var users []RequestUser
	db.DB.Model(&types.User{}).Find(&users)

	return pubsub.PubsubMessage{
		Payload: map[string]any{
			"data":   users,
			"status": "success",
			"error":  nil,
		},
		Entity:    pm.Entity,
		Operation: pm.Operation,
		Topic:     "db->auth",
		UUID:      pm.UUID,
	}, nil
}

func CreateUser(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {

	var user types.User
	err := json.Unmarshal([]byte(pm.Payload.(string)), &user)
	if err != nil {
		return pubsub.PubsubMessage{
			Payload: map[string]any{
				"data":   nil,
				"status": "Could not parse JSON",
				"error":  err.Error(),
			},
			Entity:    pm.Entity,
			Operation: pm.Operation,
			Topic:     "db->auth",
			UUID:      pm.UUID,
		}, nil
	}

	result := db.DB.Create(&user)
	if result.Error != nil {
		return pubsub.PubsubMessage{
			Payload: map[string]any{
				"data":   nil,
				"error":  result.Error.Error(),
				"status": "Could not create new user",
			},
			Entity:    pm.Entity,
			Operation: pm.Operation,
			Topic:     "db->auth",
			UUID:      pm.UUID,
		}, nil
	}

	return pubsub.PubsubMessage{
		Payload: map[string]any{
			"data":   user,
			"status": "Successfully created new user",
			"error":  nil,
		},
		Entity:    pm.Entity,
		Operation: pm.Operation,
		Topic:     "db->auth",
		UUID:      pm.UUID,
	}, nil
}

func StoreOTP(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {

	type Request struct {
		OTP		int `json:"otp"`
		Email	string	`json:"email"`
	}
	var requestBody Request
	err := json.Unmarshal([]byte(pm.Payload.(string)), &requestBody)
	if err != nil {
		return pubsub.PubsubMessage{
			Payload: map[string]any{
				"data":   nil,
				"status": "Could not parse JSON",
				"error":  err.Error(),
			},
			Entity:    pm.Entity,
			Operation: pm.Operation,
			Topic:     "db->auth",
			UUID:      pm.UUID,
		}, nil
	}

	result := db.DB.Model(&types.User{}).Where("email = ?", requestBody.Email).Update("otp", requestBody.OTP)
	if result.Error != nil {
		return pubsub.PubsubMessage{
			Payload: map[string]any{
				"error": result.Error.Error(),
			},
			Entity:    pm.Entity,
			Operation: pm.Operation,
			Topic:     "db->auth",
			UUID:      pm.UUID,
		}, nil
	}

	return pubsub.PubsubMessage{
		Payload: map[string]any{
			"data": nil,
			"status": "Successfully stored new OTP",
			"error": nil,
		},
		Entity: pm.Entity,
		Operation: pm.Operation,
		Topic: "db->auth",
		UUID: pm.UUID,
	}, nil
}

func VerifyOTP(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {

	type Request struct {
		OTP		string `json:"otp"`
		Email	string	`json:"email"`
	}
	var requestBody Request
	err := json.Unmarshal([]byte(pm.Payload.(string)), &requestBody)
	if err != nil {
		return pubsub.PubsubMessage{
			Payload: map[string]any{
				"data":   nil,
				"status": "Could not parse JSON",
				"error":  err.Error(),
			},
			Entity:    pm.Entity,
			Operation: pm.Operation,
			Topic:     "db->auth",
			UUID:      pm.UUID,
		}, nil
	}

	type UserRequest struct {
		OTP string
	}
	var user UserRequest
	result := db.DB.Model(&types.User{}).Where("email = ?", requestBody.Email).Find(&user)
	if result.Error != nil {
		return pubsub.PubsubMessage{
			Payload: map[string]any{
				"error": result.Error.Error(),
			},
			Entity:    pm.Entity,
			Operation: pm.Operation,
			Topic:     "db->auth",
			UUID:      pm.UUID,
		}, nil
	}

	if requestBody.OTP != user.OTP {
		return pubsub.PubsubMessage{
			Payload: map[string]any{
				"data": nil,
				"status": "Invalid OTP",
				"error": "Please enter correct OTP to get verified",
			},
			Entity:    pm.Entity,
			Operation: pm.Operation,
			Topic:     "db->auth",
			UUID:      pm.UUID,
		}, nil
	}

	result = db.DB.Model(&types.User{}).Where("email = ?", requestBody.Email).Update("email_verified_at", time.Now())
	if result.Error != nil {
		return pubsub.PubsubMessage{
			Payload: map[string]any{
				"error": result.Error.Error(),
			},
			Entity:    pm.Entity,
			Operation: pm.Operation,
			Topic:     "db->auth",
			UUID:      pm.UUID,
		}, nil
	}

	return pubsub.PubsubMessage{
		Payload: map[string]any{
			"data": nil,
			"status": "Successfully verified OTP",
			"error": nil,
		},
		Entity: pm.Entity,
		Operation: pm.Operation,
		Topic: "db->auth",
		UUID: pm.UUID,
	}, nil
}

func UpdateUser(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {

	var user map[string]any
	err := json.Unmarshal([]byte(pm.Payload.(string)), &user)
	if err != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"status": "Could not parse JSON",
			"error":  err.Error(),
		}, pm, "db->auth")
	}

	result := db.DB.Model(&types.User{}).Where("registration_number = ?", user["registrationNumber"]).Save(&user)
	if result.Error != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"data":   nil,
			"error":  result.Error.Error(),
			"status": "Could not update user",
		}, pm, "db->auth")
	}

	return utils.CreateRespondingPubsubMessage(map[string]any{
		"data":   user,
		"status": "Successfully updated user",
		"error":  nil,
	}, pm, "db->auth")
}
