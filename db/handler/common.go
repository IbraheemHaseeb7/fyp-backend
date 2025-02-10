package handler

import (
	"github.com/IbraheemHaseeb7/pubsub"
)

func Handle(topic string, pubsubMessage pubsub.PubsubMessage) (pubsub.PubsubMessage, error) {
	router := map[string]map[string]map[string]func(pubsub.PubsubMessage) (pubsub.PubsubMessage, error){
		"auth->db": {
			"users": {
				"READ_ALL": ReadAllUsers,
				"READ_ONE": ReadOneUser,
				"CREATE": CreateUser,
			},
			"vehicles": {
				"CREATE": CreateVehicle,
				"READ_ALL": ReadAllVehicles,
				"READ_ONE": ReadOneVehicle,
			},
		},
	}

	return router[topic][pubsubMessage.Entity][pubsubMessage.Operation](pubsubMessage)
}
