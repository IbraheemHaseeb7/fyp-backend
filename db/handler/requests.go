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

func GetAllRequests(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {
	var reqBody map[string]any
	if err := json.Unmarshal([]byte(pm.Payload.(string)), &reqBody); err != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": err.Error(),
		}, pm, "db->auth")
	}

	offset := utils.GetOffset(int(reqBody["page"].(float64)), 20)

	var requests []types.Request
	query := db.DB.Model(&types.Request{}).
		Preload("Vehicle").
		Preload("User").
		Limit(20).
		Offset(offset).
		Where("status = ?", reqBody["status"])
	var result *gorm.DB

	if reqBody["me"] != "false" {
		result = query.Where("user_id = ?", reqBody["me"]).Find(&requests)
	}

	result = query.Find(&requests)
	if result.Error != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": result.Error.Error(),
		}, pm, "db->auth")
	}

	return utils.CreateRespondingPubsubMessage(map[string]any{
		"data": requests,
	}, pm, "db->auth")
}

func GetSingleRequest(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {

	var reqBody types.Request
	if err := json.Unmarshal([]byte(pm.Payload.(string)), &reqBody); err != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": err.Error(),
		}, pm, "db->auth")
	}

	var request types.Request
	result := db.DB.Model(&types.Request{}).Preload("Vehicle").
		Preload("User", func (db *gorm.DB) *gorm.DB {
			return db.Select("id, name, email, registration_number")
		}).
		Where("id = ?", reqBody.ID).First(&request)
	if result.Error != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": result.Error.Error(),
		}, pm, "db->auth")
	}

	return utils.CreateRespondingPubsubMessage(map[string]any{
		"data": request,
	}, pm, "db->auth")
}

func CreateRequest(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {

	var reqBody types.Request
	if err := json.Unmarshal([]byte(pm.Payload.(string)), &reqBody); err != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": err.Error(),
		}, pm, "db->auth")
	}

	// checking for active requests under this user_id
	var count int64
	result := db.DB.Model(&types.Request{}).
		Where("user_id = ? AND status <> ? AND status <> ? AND status <> ?", reqBody.UserID, "completed", "expired", "rejected").
		Count(&count)
	if result.Error != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": result.Error.Error(),
		}, pm, "db->auth")
	}
	fmt.Println("Count: ", count, " UserID: ", reqBody.UserID)
	if count > 0 {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": "You already have a an active request",
			"status": "Cannot create new ride request",
		}, pm, "db->auth")
	}

	// checking if the role is rider, then they must have 1 or more vehicles added
	if reqBody.OriginatorRole == "rider" {
		result := db.DB.Model(&types.Vehicle{}).
			Where("user_id = ?", reqBody.UserID).
			Count(&count)
		if result.Error != nil {
			return utils.CreateRespondingPubsubMessage(map[string]any{
				"error": result.Error.Error(),
			}, pm, "db->auth")
		}
		if count == 0  {
			return utils.CreateRespondingPubsubMessage(map[string]any{
				"error": "You don't have any vehicles added to generate request",
				"status": "Cannot create new ride request",
			}, pm, "db->auth")
		}
	}

	result = db.DB.Create(&reqBody)
	if result.Error != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": result.Error.Error(),
		}, pm, "db->auth")
	}

	return utils.CreateRespondingPubsubMessage(map[string]any{
		"data": reqBody,
	}, pm, "db->auth")
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
		}, pm, "db->auth")
	}

	result := db.DB.Where("id = ? and user_id = ?", reqBody.ID, reqBody.UserID).Delete(&types.Request{})
	if result.Error != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": result.Error.Error(),
		}, pm, "db->auth")
	}

	db.DB.Where("request_id = ?", reqBody.ID).Delete(&types.Request{})

	return utils.CreateRespondingPubsubMessage(map[string]any{
		"data": fmt.Sprintf("Deleted rows count: %d", result.RowsAffected),
	}, pm, "db->auth")
}

func SetStatus(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {

	var reqBody map[string]any
	if err := json.Unmarshal([]byte(pm.Payload.(string)), &reqBody); err != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": err.Error(),
		}, pm, "db->auth")
	}

	result := db.DB.Model(&types.Request{}).Where("id = ?", reqBody["id"]).Update("status", reqBody["status"])
	if result.Error != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": result.Error.Error(),
		}, pm, "db->auth")
	}

	return utils.CreateRespondingPubsubMessage(map[string]any{
		"data": fmt.Sprintf("Successfully updated request's status"),
	}, pm, "db->auth")
}

func  GetMyProposalForARequest(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {
	var reqBody map[string]any
	if err := json.Unmarshal([]byte(pm.Payload.(string)), &reqBody); err != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": err.Error(),
		}, pm, "db->auth")
	}

	var data types.Request
	result := db.DB.Model(&types.Request{}).
		Where("request_id = ? AND user_id = ?", reqBody["request_id"], reqBody["user_id"]).
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

	if result.RowsAffected == 0 {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": "Not found",
		}, pm, "db->auth")
	}

	return utils.CreateRespondingPubsubMessage(map[string]any{
		"data": data,
	}, pm, "db->auth")
}
