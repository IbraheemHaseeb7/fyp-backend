package handler

import (
	"encoding/json"
	"fmt"

	"github.com/IbraheemHaseeb7/fyp-backend/db"
	"github.com/IbraheemHaseeb7/fyp-backend/utils"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/IbraheemHaseeb7/types"
)

func GetAllRequests(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {
	var reqBody map[string]any
	if err := json.Unmarshal([]byte(pm.Payload.(string)), &reqBody); err != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": err.Error(),
		}, pm, "db->img")
	}

	offset := utils.GetOffset(int(reqBody["page"].(float64)), 20)

	var requests []types.Request
	result := db.DB.Model(&types.Request{}).Limit(20).Offset(offset).Where("status = ?", reqBody["status"]).Find(&requests)
	if result.Error != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": result.Error.Error(),
		}, pm, "db->img")
	}

	return utils.CreateRespondingPubsubMessage(map[string]any{
		"data": requests,
	}, pm, "db->img")
}

func GetSingleRequest(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {

	var reqBody types.Request
	if err := json.Unmarshal([]byte(pm.Payload.(string)), &reqBody); err != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": err.Error(),
		}, pm, "db->img")
	}

	var request types.Request
	result := db.DB.Model(&types.Request{}).Where("id = ?", reqBody.ID).First(&request)
	if result.Error != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": result.Error.Error(),
		}, pm, "db->img")
	}

	return utils.CreateRespondingPubsubMessage(map[string]any{
		"data": request,
	}, pm, "db->img")
}

func CreateRequest(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {

	var reqBody types.Request
	if err := json.Unmarshal([]byte(pm.Payload.(string)), &reqBody); err != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": err.Error(),
		}, pm, "db->img")
	}

	result := db.DB.Create(&reqBody)
	if result.Error != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": result.Error.Error(),
		}, pm, "db->img")
	}

	return utils.CreateRespondingPubsubMessage(map[string]any{
		"data": reqBody,
	}, pm, "db->img")
}

func UpdateRequest(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {
	var request map[string]any
	err := json.Unmarshal([]byte(pm.Payload.(string)), &request)
	if err != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"status": "Could not parse JSON",
			"error":  err.Error(),
		}, pm, "db->auth")
	}

	result := db.DB.Model(&types.Request{}).Where("id = ? and user_id = ?", request["id"], request["user_id"]).Save(&request)
	if result.Error != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"data":   nil,
			"error":  result.Error.Error(),
			"status": "Could not update request",
		}, pm, "db->auth")
	}

	return utils.CreateRespondingPubsubMessage(map[string]any{
		"data":		request,
		"status": 	"Successfully updated request",
		"error":  	nil,
	}, pm, "db->auth")
}

func DeleteRequest(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {
	var reqBody types.Request
	if err := json.Unmarshal([]byte(pm.Payload.(string)), &reqBody); err != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": err.Error(),
		}, pm, "db->img")
	}

	result := db.DB.Where("id = ? and user_id = ?", reqBody.ID, reqBody.UserID).Delete(&types.Request{})
	if result.Error != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": result.Error.Error(),
		}, pm, "db->img")
	}

	return utils.CreateRespondingPubsubMessage(map[string]any{
		"data": fmt.Sprintf("Deleted rows count: %d", result.RowsAffected),
	}, pm, "db->img")
}
