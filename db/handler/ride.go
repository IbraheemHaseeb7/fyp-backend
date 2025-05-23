package handler

import (
	"encoding/json"
	"fmt"

	"github.com/IbraheemHaseeb7/fyp-backend/db"
	"github.com/IbraheemHaseeb7/fyp-backend/utils"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/IbraheemHaseeb7/types"
	"gorm.io/gorm"
)

func GetAllRides(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {
	type Query struct {
		UserID string `json:"user_id"`
		PageNo int    `json:"page"`
	}
	var query Query
	if err := json.Unmarshal([]byte(pm.Payload.(string)), &query); err != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": err.Error(),
		}, pm, "db->auth")
	}

	offset := utils.GetOffset(query.PageNo, 20)
	var rides []types.Ride
	result := db.DB.Model(&types.Ride{}).Limit(20).Offset(offset).Find(&rides)
	if result.Error != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": result.Error.Error(),
		}, pm, "db->auth")
	}

	return utils.CreateRespondingPubsubMessage(map[string]any{
		"data": rides,
	}, pm, "db->auth")
}

func GetSingleRide(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {
	type Query struct {
		ID     string `json:"id"`
		UserID string `json:"user_id"`
	}

	var query Query
	if err := json.Unmarshal([]byte(pm.Payload.(string)), &query); err != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": err.Error(),
		}, pm, "db->auth")
	}

	var ride types.Ride
	result := db.DB.Model(&types.Ride{}).Where("id = ? AND user_id = ?", query.ID, query.UserID).First(&ride)
	if result.Error != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": result.Error.Error(),
		}, pm, "db->auth")
	}
	if result.RowsAffected == 0 {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": "Ride not found",
		}, pm, "db->auth")
	}

	return utils.CreateRespondingPubsubMessage(map[string]any{
		"data": ride,
	}, pm, "db->auth")
}

func CreateRide(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {
	var requestBody map[string]any
	if err := json.Unmarshal([]byte(pm.Payload.(string)), &requestBody); err != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": err.Error(),
		}, pm, "db->auth")
	}
	result := db.DB.Model(&types.Ride{}).Create(&requestBody)
	if result.Error != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": result.Error.Error(),
		}, pm, "db->auth")
	}

	return utils.CreateRespondingPubsubMessage(map[string]any{
		"message": "Successfully created ride",
	}, pm, "db->auth")
}

func UpdateRide(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {
	var requestBody map[string]any
	if err := json.Unmarshal([]byte(pm.Payload.(string)), &requestBody); err != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": err.Error(),
		}, pm, "db->auth")
	}

	result := db.DB.Model(&types.Ride{}).Where("id = ? AND user_id = ?", requestBody["id"], requestBody["user_id"]).Updates(requestBody)
	if result.Error != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": result.Error.Error(),
		}, pm, "db->auth")
	}
	if result.RowsAffected == 0 {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": "Ride not found",
		}, pm, "db->auth")
	}

	return utils.CreateRespondingPubsubMessage(map[string]any{
		"message": "Successfully updated ride",
	}, pm, "db->auth")
}

func DeleteRide(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {
	type Query struct {
		ID     string `json:"id"`
		UserID string `json:"user_id"`
	}
	var query Query
	if err := json.Unmarshal([]byte(pm.Payload.(string)), &query); err != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": err.Error(),
		}, pm, "db->auth")
	}

	result := db.DB.Model(&types.Ride{}).Where("id = ? AND user_id = ?", query.ID, query.UserID).Delete(&types.Ride{})
	if result.Error != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": result.Error.Error(),
		}, pm, "db->auth")
	}
	if result.RowsAffected == 0 {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": "Ride not found",
		}, pm, "db->auth")
	}

	return utils.CreateRespondingPubsubMessage(map[string]any{
		"message": "Successfully deleted ride",
	}, pm, "db->auth")
}


func ActiveRide(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {
	fmt.Println("ActiveRide")
	type Query struct {
		UserID float64 `json:"user_id"`
	}
	var query Query
	if err := json.Unmarshal([]byte(pm.Payload.(string)), &query); err != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": err.Error(),
		}, pm, "db->auth")
	}

	var ride types.Ride
	result := db.DB.Model(&types.Ride{}).
		Where("request_id = (select id from requests where user_id = ? and status = ?)", query.UserID, "matched").
		Preload("Driver", func (db *gorm.DB) *gorm.DB {
			return db.Select("id, name, email")
		}).
		Preload("Passenger", func (db *gorm.DB) *gorm.DB {
			return db.Select("id, name, email")
		}).
		Preload("Vehicle").
		Preload("Request").
		Preload("Proposal").
		First(&ride)
	if result.Error != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": result.Error.Error(),
		}, pm, "db->auth")
	}
	if result.RowsAffected == 0 {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": "Ride not found",
		}, pm, "db->auth")
	}

	return utils.CreateRespondingPubsubMessage(map[string]any{
		"message": "Successfully fetched ride",
		"data":    ride,
	}, pm, "db->auth")
}
