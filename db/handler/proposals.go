package handler

import (
	"encoding/json"

	"github.com/IbraheemHaseeb7/fyp-backend/db"
	"github.com/IbraheemHaseeb7/fyp-backend/utils"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/IbraheemHaseeb7/types"
	"gorm.io/gorm"
)

func GetAllProposals(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {
	var reqBody map[string]any
	err := json.Unmarshal([]byte(pm.Payload.(string)), &reqBody)
	if err != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": err.Error(),
		}, pm, "db->auth")
	}

	var temp map[string]any
	result := db.DB.Model(&types.Request{}).Where("id = ? AND user_id = ?", reqBody["request_id"], reqBody["user_id"]).Find(&temp)
	if result.RowsAffected == 0 {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": "Not authorised to view these proposals",
		}, pm, "db->auth")
	}

	var data []types.Request
	result = db.DB.Model(&types.Request{}).
		Where("request_id = ?", reqBody["request_id"]).
		Preload("User", func (db *gorm.DB) *gorm.DB {
			return db.Select("id, name, email, registration_number")
		}).
		Preload("Vehicle", func (db *gorm.DB) *gorm.DB {
			return db.Select("*")
		}).
		Find(&data)
	if result.Error != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": result.Error.Error(),
		}, pm, "db->auth")
	}

	return utils.CreateRespondingPubsubMessage(map[string]any{
		"data": data,
	}, pm, "db->auth")
}

func GetAllMyProposals(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {
	var reqBody map[string]any
	err := json.Unmarshal([]byte(pm.Payload.(string)), &reqBody)
	if err != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": err.Error(),
		}, pm, "db->auth")
	}

	var data []types.Request
	result := db.DB.Model(&types.Request{}).
		Preload("User", func (db *gorm.DB) *gorm.DB {
			return db.Select("id, name, email, registration_number")
		}).
		Preload("Vehicle", func (db *gorm.DB) *gorm.DB {
			return db.Select("*")
		}).
		Where("user_id = ? AND status = ?", reqBody["user_id"],  "proposal").Find(&data)
	if result.Error != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": result.Error.Error(),
		}, pm, "db->auth")
	}

	return utils.CreateRespondingPubsubMessage(map[string]any{
		"data": data,
	}, pm, "db->auth")
}
