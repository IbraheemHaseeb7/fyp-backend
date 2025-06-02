package handler

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/IbraheemHaseeb7/fyp-backend/db"
	"github.com/IbraheemHaseeb7/fyp-backend/utils"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/IbraheemHaseeb7/types"
)


func CreateChatRoom(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {

	// validating and binding request
	type Request struct {
		RequestID string `json:"request_id"`
		ProposalID string `json:"proposal_id"`
	}
	var requestBody Request
	if err := json.Unmarshal([]byte(pm.Payload.(string)), &requestBody); err != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": err.Error(),
		}, pm, "db->auth")
	}

	// convert string to int64
	requestID, err := strconv.ParseInt(requestBody.RequestID, 10, 64)
	if err != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": "Invalid request ID",
		}, pm, "db->auth")
	}
	proposalID, err := strconv.ParseInt(requestBody.ProposalID, 10, 64)
	if err != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": "Invalid proposal ID",
		}, pm, "db->auth")
	}

	// reflecting changes in the database
	room := types.Room{
		RequestID: requestID,
		ProposalID: proposalID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	result := db.DB.Model(&types.Room{}).Create(&room)

	if result.Error != nil {
		// if there is duplicate entry error then
		if strings.Split(result.Error.Error(), " ")[2] == "(23000):" {

			// fetch the room from the database
			result = db.DB.Model(&types.Room{}).Where("request_id = ? AND proposal_id = ?", requestID, proposalID).First(&room)
			if result.Error != nil {
				return utils.CreateRespondingPubsubMessage(map[string]any{
					"error": result.Error.Error(),
				}, pm, "db->auth")
			}

			return utils.CreateRespondingPubsubMessage(map[string]any{
				"data": room,
			}, pm, "db->auth")
		}
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": result.Error.Error(),
		}, pm, "db->auth")
	}

	// sending response
	return utils.CreateRespondingPubsubMessage(map[string]any{
		"data": room,
	}, pm, "db->auth")
}

func GetChatMessages(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {
	// validating and binding request
	type Request struct {
		RoomID string `json:"room_id"`
		Page   string `json:"page"`
		UserID int64  `json:"user_id"`
	}
	var requestBody Request
	if err := json.Unmarshal([]byte(pm.Payload.(string)), &requestBody); err != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": err.Error(),
		}, pm, "db->auth")
	}

	roomID, err := strconv.ParseInt(requestBody.RoomID, 10, 64)
	if err != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": "Invalid room ID",
		}, pm, "db->auth")
	}

	page, err := strconv.Atoi(requestBody.Page)
	if err != nil || page < 1 {
		page = 1 // Default to page 1 if not provided or invalid
	}

	type Message struct {
		ID        int64     `json:"id"`
		RoomID    int64     `json:"room_id"`
		UserID    int64     `json:"user_id"`
		Message   string    `json:"message"`	
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	var messages []Message
	result := db.DB.Model(&Message{}).Where("room_id = ? and room_id in (select id from rooms where rooms.request_id in (select id from requests where user_id = ?) or rooms.proposal_id in (select id from requests where user_id = ?))", roomID, requestBody.UserID, requestBody.UserID).Offset((page - 1) * 20).Limit(20).Order("created_at DESC").Find(&messages)
	if result.Error != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": result.Error.Error(),
		}, pm, "db->auth")
	}

	return utils.CreateRespondingPubsubMessage(map[string]any{
		"data": messages,
	}, pm, "db->auth")
}
