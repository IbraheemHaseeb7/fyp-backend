package handler

import (
	"encoding/json"
	"time"

	"github.com/IbraheemHaseeb7/fyp-backend/db"
	"github.com/IbraheemHaseeb7/fyp-backend/utils"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/IbraheemHaseeb7/types"
)

func AddFeedback(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {
	var feedback map[string]any
	if err := json.Unmarshal([]byte(pm.Payload.(string)), &feedback); err != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": err.Error(),
		}, pm, "db->auth")
	}
	feedback["created_at"] = time.Now().Format(time.RFC3339)

	result := db.DB.Model(&types.Feedback{}).Create(&feedback)
	if result.Error != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": result.Error.Error(),
		}, pm, "db->auth")
	}

	return utils.CreateRespondingPubsubMessage(map[string]any{
		"message": "Successfully added feedback",
	}, pm, "db->auth")
}
