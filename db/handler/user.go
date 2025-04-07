package handler

import (
	"encoding/json"

	"github.com/IbraheemHaseeb7/fyp-backend/db"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/IbraheemHaseeb7/types"
)

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

	var user types.User
	db.DB.Where("registration_number = ?", query.RegistrationNumber).First(&user)

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

	var users []types.User
	db.DB.Find(&users)

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
