package handler

import (
	"github.com/IbraheemHaseeb7/fyp-backend/auth"
	"github.com/IbraheemHaseeb7/pubsub"
)

func VerifyToken(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {
	payload, ok := pm.Payload.(map[string]any)
	if !ok {
		return pubsub.PubsubMessage{
			Payload: map[string]any{
				"verified": false,
			},
			Entity:    pm.Entity,
			Operation: pm.Operation,
			Topic:     "auth->img",
			UUID:      pm.UUID,
		}, nil
	}

	token, ok := payload["token"].(string)
	if !ok {
		return pubsub.PubsubMessage{
			Payload: map[string]any{
				"verified": false,
			},
			Entity:    pm.Entity,
			Operation: pm.Operation,
			Topic:     "auth->img",
			UUID:      pm.UUID,
		}, nil
	}

	err := auth.VerifyToken(token)
	if err != nil {
		return pubsub.PubsubMessage{
			Payload: map[string]any{
				"verified": false,
			},
			Entity:    pm.Entity,
			Operation: pm.Operation,
			Topic:     "auth->img",
			UUID:      pm.UUID,
		}, nil
	}

	return pubsub.PubsubMessage{
		Payload: map[string]any{
			"verified": true,
		},
		Entity:    pm.Entity,
		Operation: pm.Operation,
		Topic:     "auth->img",
		UUID:      pm.UUID,
	}, nil
}

func GetClaims(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {
	payload, ok := pm.Payload.(map[string]any)
	if !ok {
		return pubsub.PubsubMessage{
			Payload: map[string]any{
				"verified": false,
			},
			Entity:    pm.Entity,
			Operation: pm.Operation,
			Topic:     "auth->img",
			UUID:      pm.UUID,
		}, nil
	}

	token, ok := payload["token"].(string)
	if !ok {
		return pubsub.PubsubMessage{
			Payload: map[string]any{
				"verified": false,
			},
			Entity:    pm.Entity,
			Operation: pm.Operation,
			Topic:     "auth->img",
			UUID:      pm.UUID,
		}, nil
	}
	claims, err := auth.GetClaimsFromToken(token)
	if err != nil {
		return pubsub.PubsubMessage{
			Payload: map[string]any{
				"verified": false,
				"error": err.Error(),
			},
			Entity:    pm.Entity,
			Operation: pm.Operation,
			Topic:     "auth->img",
			UUID:      pm.UUID,
		}, nil
	}

	return pubsub.PubsubMessage{
		Payload: claims,
		Entity:    pm.Entity,
		Operation: pm.Operation,
		Topic:     "auth->img",
		UUID:      pm.UUID,
	}, nil
}
