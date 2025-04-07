package handler

import (
	"github.com/IbraheemHaseeb7/pubsub"
)

func VerifyToken(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {
	return pm, nil
}

func GetClaims(pm pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {
	return pm, nil
}
