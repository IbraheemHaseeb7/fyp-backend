package handler

import (
	"github.com/IbraheemHaseeb7/fyp-backend/utils"
	"github.com/IbraheemHaseeb7/pubsub"
)

func GetAllRequests(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {
	return utils.CreateRespondingPubsubMessage(map[string]any{
		"data": "Successfully created new request",
	}, pm, "db->img")
}

func GetSingleRequest(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {
	return utils.CreateRespondingPubsubMessage(map[string]any{
		"data": "Successfully created new request",
	}, pm, "db->img")
}

func CreateRequest(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {
	return utils.CreateRespondingPubsubMessage(map[string]any{
		"data": "Successfully created new request",
	}, pm, "db->img")
}

func UpdateRequest(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {
	return utils.CreateRespondingPubsubMessage(map[string]any{
		"data": "Successfully created new request",
	}, pm, "db->img")
}

func DeleteRequest(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {
	return utils.CreateRespondingPubsubMessage(map[string]any{
		"data": "Successfully created new request",
	}, pm, "db->img")
}
