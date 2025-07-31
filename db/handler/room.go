package handler

import (
	"encoding/json"

	"github.com/IbraheemHaseeb7/fyp-backend/db"
	"github.com/IbraheemHaseeb7/fyp-backend/utils"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/IbraheemHaseeb7/types"
	"gorm.io/gorm"
)

func GetRoom(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {
	var reqBody map[string]any
	if err := json.Unmarshal([]byte(pm.Payload.(string)), &reqBody); err != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": err.Error(),
		}, pm, "db->auth")
	}

	var room types.Room
	result := db.DB.Model(&types.Room{}).
		Preload("Request").
		Preload("Proposal").
		Preload("Proposal.User", func (db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "email", "device_token", "profile_uri")
		}).
		Preload("Request.User", func (db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "email", "device_token", "profile_uri")
		}).
		Where("id = ?", reqBody["room_id"]).First(&room)
	if result.Error != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": "Room not found",
		}, pm, "db->auth")
	}


	return utils.CreateRespondingPubsubMessage(map[string]any{
		"data": room,
	}, pm, "db->auth")
}

func GetRoomWithIDs(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {
	var reqBody map[string]any
	if err := json.Unmarshal([]byte(pm.Payload.(string)), &reqBody); err != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": err.Error(),
		}, pm, "db->auth")
	}

	var room types.Room
	result := db.DB.Model(&types.Room{}).
		Preload("Request").
		Preload("Proposal").
		Preload("Proposal.User", func (db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "email", "device_token", "profile_uri")
		}).
		Preload("Request.User", func (db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "email", "device_token", "profile_uri")
		}).
		Where("request_id = ? AND proposal_id = ?", reqBody["request_id"], reqBody["proposal_id"]).First(&room)
	if result.Error != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": "Room not found",
		}, pm, "db->auth")
	}

	return utils.CreateRespondingPubsubMessage(map[string]any{
		"data": room,
	}, pm, "db->auth")
}
