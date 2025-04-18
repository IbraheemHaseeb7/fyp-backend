package handler

import (
	"encoding/json"
	"time"

	"github.com/IbraheemHaseeb7/fyp-backend/db"
	"github.com/IbraheemHaseeb7/fyp-backend/utils"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/IbraheemHaseeb7/types"
)

func VerifyCard(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {

	// validating and binding request
	type Request struct {
		RegistrationNumber string `json:"registrationNumber"`
		CardStatus	string	`json:"cardStatus"`
	}
	var requestBody Request
	if err := json.Unmarshal([]byte(pm.Payload.(string)), &requestBody); err != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": err.Error(),
		}, pm, "db->img")
	}
	if requestBody.CardStatus == "false" {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": "card could not be verified",
		}, pm, "db->img")
	} 

	// reflecting changes in the database
	result := db.DB.Model(&types.User{}).
		Where("registration_number = ?", requestBody.RegistrationNumber).
		Update("card_verified_at", time.Now())

	if result.Error != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": result.Error.Error(),
		}, pm, "db->img")
	}

	// sending response
	return utils.CreateRespondingPubsubMessage(map[string]any{
		"data": requestBody,
	}, pm, "db->img")
}

func VerifySelfie(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {

	// validating and binding request
	type Request struct {
		RegistrationNumber string `json:"registrationNumber"`
		CardStatus	string	`json:"cardStatus"`
	}
	var requestBody Request
	if err := json.Unmarshal([]byte(pm.Payload.(string)), &requestBody); err != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": err.Error(),
		}, pm, "db->img")
	}
	if requestBody.CardStatus == "false" {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": "card could not be verified",
		}, pm, "db->img")
	} 

	// reflecting changes in the database
	result := db.DB.Model(&types.User{}).
		Where("registration_number = ?", requestBody.RegistrationNumber).
		Update("selfie_verified_at", time.Now())

	if result.Error != nil {
		return utils.CreateRespondingPubsubMessage(map[string]any{
			"error": result.Error.Error(),
		}, pm, "db->img")
	}

	// sending response
	return utils.CreateRespondingPubsubMessage(map[string]any{
		"data": requestBody,
	}, pm, "db->img")
}
