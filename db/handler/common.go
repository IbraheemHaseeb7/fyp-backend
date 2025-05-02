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
				"LOGIN":	Login,
				"CREATE":   CreateUser,
				"STORE_OTP": StoreOTP,
				"VERIFY_OTP": VerifyOTP,
				"UPDATE_ONE": UpdateUser,
			},
			"vehicles": {
				"CREATE":   CreateVehicle,
				"UPDATE":   UpdateVehicle,
				"DELETE":   DeleteVehicle,
				"READ_ALL": ReadAllVehicles,
				"READ_ONE": ReadOneVehicle,
			},
			"requests": {
				"CREATE":	CreateRequest,
				"UPDATE":	UpdateRequest,
				"DELETE":	DeleteRequest,
				"READ_ALL":	GetAllRequests,
				"READ_ONE": GetSingleRequest,
				"SET_STATUS":SetStatus,
			},
			"proposals": {
				"GET_ALL":	GetAllProposals,
				"GET_ALL_MY":	GetAllMyProposals,
			},
		},
		"img->db": {
			"files": {
				"VERIFY_CARD": VerifyCard,
				"VERIFY_SELFIE": VerifySelfie,
			},
		},
	}

	return router[topic][pubsubMessage.Entity][pubsubMessage.Operation](pubsubMessage)
}
