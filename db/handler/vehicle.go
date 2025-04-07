package handler

import (
	"encoding/json"

	"github.com/IbraheemHaseeb7/fyp-backend/db"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/IbraheemHaseeb7/types"
)

func ReadAllVehicles(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {

	type Query struct {
		RegNo string `json:"regNo"`
	}
	var query Query
	if err := json.Unmarshal([]byte(pm.Payload.(string)), &query); err != nil {
		return pubsub.PubsubMessage{
			Payload: map[string]any{
				"error": err.Error(),
			},
			Entity:    pm.Entity,
			Operation: pm.Operation,
			Topic:     "db->auth",
			UUID:      pm.UUID,
		}, nil
	}

	var vehicles []types.Vehicle
	result := db.DB.Raw("SELECT * FROM vehicles WHERE user_id = (SELECT id FROM users WHERE registration_number = ?)", query.RegNo).Scan(&vehicles)

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
			"data": vehicles,
		},
		Entity:    pm.Entity,
		Operation: pm.Operation,
		Topic:     "db->auth",
		UUID:      pm.UUID,
	}, nil
}

func ReadOneVehicle(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {

	type Query struct {
		ID    int64
		RegNo string
	}
	var query Query
	err := json.Unmarshal([]byte(pm.Payload.(string)), &query)
	if err != nil {
		return pubsub.PubsubMessage{
			Payload: map[string]any{
				"error": err.Error(),
			},
			Entity:    pm.Entity,
			Operation: pm.Operation,
			Topic:     "db->auth",
			UUID:      pm.UUID,
		}, nil
	}

	var vehicle types.Vehicle
	result := db.DB.Raw("SELECT * FROM vehicles WHERE id = ? AND user_id = (SELECT id FROM users WHERE registration_number = ?)", query.ID, query.RegNo).Scan(&vehicle)

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
			"data": vehicle,
		},
		Entity:    pm.Entity,
		Operation: pm.Operation,
		Topic:     "db->auth",
		UUID:      pm.UUID,
	}, nil
}

func CreateVehicle(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {

	var vehicle types.Vehicle
	err := json.Unmarshal([]byte(pm.Payload.(string)), &vehicle)
	if err != nil {
		return pubsub.PubsubMessage{
			Payload: map[string]any{
				"error": err.Error(),
			},
			Entity:    pm.Entity,
			Operation: pm.Operation,
			Topic:     "db->auth",
			UUID:      pm.UUID,
		}, nil
	}

	result := db.DB.Create(&vehicle)
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
		Payload:   map[string]any{},
		Entity:    pm.Entity,
		Operation: pm.Operation,
		Topic:     "db->auth",
		UUID:      pm.UUID,
	}, nil
}

func UpdateVehicle(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {

	var vehicle types.Vehicle
	err := json.Unmarshal([]byte(pm.Payload.(string)), &vehicle)
	if err != nil {
		return pubsub.PubsubMessage{
			Payload: map[string]any{
				"error": err.Error(),
			},
			Entity:    pm.Entity,
			Operation: pm.Operation,
			Topic:     "db->auth",
			UUID:      pm.UUID,
		}, nil
	}

	result := db.DB.Save(&vehicle)
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
		Payload:   map[string]any{},
		Entity:    pm.Entity,
		Operation: pm.Operation,
		Topic:     "db->auth",
		UUID:      pm.UUID,
	}, nil
}

func DeleteVehicle(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {

	type Query struct {
		ID int64 `json:"id"`
		VehicleID int64 `json:"vId"`
	}
	var query Query
	err := json.Unmarshal([]byte(pm.Payload.(string)), &query)
	if err != nil {
		return pubsub.PubsubMessage{
			Payload: map[string]any{
				"error": err.Error(),
			},
			Entity:    pm.Entity,
			Operation: pm.Operation,
			Topic:     "db->auth",
			UUID:      pm.UUID,
		}, nil
	}

	result := db.DB.Where("user_id = ? and id = ?", query.ID, query.VehicleID).Delete(&types.Vehicle{})
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
		Payload:   map[string]any{},
		Entity:    pm.Entity,
		Operation: pm.Operation,
		Topic:     "db->auth",
		UUID:      pm.UUID,
	}, nil
}
