package handler

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/IbraheemHaseeb7/fyp-backend/db"
	"github.com/IbraheemHaseeb7/fyp-backend/utils"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/IbraheemHaseeb7/types"
)

func SendMessage(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {
	// validating and binding request
	type Request struct {
		RoomID string `json:"room_id"`
		Message string `json:"message"`
		Sender string `json:"sender"`
	}
	var requestBody Request
	if err := json.Unmarshal([]byte(pm.Payload.(string)), &requestBody); err != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": err.Error(),
		}, pm, "db->auth")
	}

	// convert string to int64
	roomID, err := strconv.ParseInt(requestBody.RoomID, 10, 64)
	if err != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": "Invalid room ID",
		}, pm, "db->auth")
	}
	userID, err := strconv.ParseInt(requestBody.Sender, 10, 64)
	if err != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": "Invalid user ID",
		}, pm, "db->auth")
	}

	// reflecting changes in the database
	room := types.Message{
		Message: requestBody.Message,
		RoomID: roomID,
		UserID: userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result := db.DB.Model(&types.Message{}).Create(&room)

	if result.Error != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": result.Error.Error(),
		}, pm, "db->auth")
	}

	// sending response
	return utils.CreateRespondingPubsubMessage(map[string]any{
		"data": room,
	}, pm, "db->auth")
}

