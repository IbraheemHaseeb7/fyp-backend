package handler

import (
	"encoding/json"

	"github.com/IbraheemHaseeb7/fyp-backend/db"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/IbraheemHaseeb7/types"
)

func ReadAllVehicles(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {
	var vehicles []types.Vehicle
	result := db.DB.Find(&vehicles)
	if result.Error != nil {
		return pubsub.PubsubMessage{
			Payload: map[string]any{
				"error": result.Error.Error(),
			},
			Entity: pm.Entity,
			Operation: pm.Operation,
			Topic: "db->auth",
			UUID: pm.UUID,
		}, nil
	}

	return pubsub.PubsubMessage{
		Payload: map[string]any{
			"data": vehicles,
		},
		Entity: pm.Entity,
		Operation: pm.Operation,
		Topic: "db->auth",
		UUID: pm.UUID,
	}, nil
}

func ReadOneVehicle(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {

	type Query struct {
		ID int64
	}
	var query Query
	err := json.Unmarshal([]byte(pm.Payload.(string)), &query)
	if err != nil {
		return pubsub.PubsubMessage{
			Payload: map[string]any{
				"error": err.Error(),
			},
			Entity: pm.Entity,
			Operation: pm.Operation,
			Topic: "db->auth",
			UUID: pm.UUID,
		}, nil
	}

	var vehicle types.Vehicle
	result := db.DB.Where("id = ?", query.ID).First(&vehicle)
	if result.Error != nil {
		return pubsub.PubsubMessage{
			Payload: map[string]any{
				"error": result.Error.Error(),
			},
			Entity: pm.Entity,
			Operation: pm.Operation,
			Topic: "db->auth",
			UUID: pm.UUID,
		}, nil
	}

	return pubsub.PubsubMessage{
		Payload: map[string]any{
			"data": vehicle,
		},
		Entity: pm.Entity,
		Operation: pm.Operation,
		Topic: "db->auth",
		UUID: pm.UUID,
	}, nil
}

func CreateVehicle(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {

	vehicle := types.Vehicle{
		Make: "Yamaha",
		Model: "YBR",
		Year: 2023,
		Type: "bike",
		UserID: 133,
	}
	result := db.DB.Create(&vehicle)
	if result.Error != nil {
		return pubsub.PubsubMessage{
			Payload: map[string]any{
				"error": result.Error.Error(),
			},
			Entity: pm.Entity,
			Operation: pm.Operation,
			Topic: "db->auth",
			UUID: pm.UUID,
		}, nil
	}

	return pubsub.PubsubMessage{
		Payload: map[string]any{
		},
		Entity: pm.Entity,
		Operation: pm.Operation,
		Topic: "db->auth",
		UUID: pm.UUID,
	}, nil
}
