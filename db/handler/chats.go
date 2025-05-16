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
