package handler

import (
	"github.com/IbraheemHaseeb7/pubsub"
)

func Handle(topic string, pubsubMessage pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {
	router := map[string]map[string]map[string]func(pubsub.PubsubMessage) (pubsub.PubsubMessage, error){
		"db->auth": {
			"users": {
				"READ_ALL": ReadAllUsers,
				"READ_ONE": ReadOneUser,
				"CREATE":   CreateUser,
			},
			"vehicles": {
				"READ_ALL": ReadAllVehicles,
				"READ_ONE": ReadOneVehicle,
				"CREATE":   CreateVehicle,
			},
		},
	}

	return router[topic][pubsubMessage.Entity][pubsubMessage.Operation](pubsubMessage)
}
