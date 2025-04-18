package utils

import "github.com/IbraheemHaseeb7/pubsub"

func GetOffset(page, pageSize int) int {
	return (page - 1) * pageSize
}

func CreateRespondingPubsubMessage(payload map[string]any, pm pubsub.PubsubMessage, topic string) (pubsub.PubsubMessage, error) {
	return pubsub.PubsubMessage{
		Payload: payload,
		Entity:    pm.Entity,
		Operation: pm.Operation,
		Topic:     topic,
		UUID:      pm.UUID,
	}, nil
}
