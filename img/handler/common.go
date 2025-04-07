package handler

import (
	"github.com/IbraheemHaseeb7/pubsub"
)

func Handle(topic string, pubsubMessage pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {
	router := map[string]map[string]map[string]func(pubsub.PubsubMessage) (pubsub.PubsubMessage, error){
		"auth->img": {
			"files": {
				"VERIFY_TOKEN": VerifyToken,
				"GET_CLAIMS":   GetClaims,
			},
		},
	}

	return router[topic][pubsubMessage.Entity][pubsubMessage.Operation](pubsubMessage)
}
