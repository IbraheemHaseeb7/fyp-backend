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
				"LOGIN": 	ReadOneUser,
				"CREATE":   CreateUser,
				"STORE_OTP": StoreOTP,
				"VERIFY_OTP": VerifyOTP,
				"UPDATE_ONE": UpdateUser,
			},
			"vehicles": {
				"READ_ALL": ReadAllVehicles,
				"READ_ONE": ReadOneVehicle,
				"CREATE":   CreateVehicle,
				"UPDATE":   UpdateVehicle,
				"DELETE":   DeleteVehicle,
			},
			"requests": {
				"READ_ALL": GetAllRequests,
				"READ_ONE": GetSingleRequest,
				"CREATE":	CreateRequest,
				"UPDATE":	UpdateRequest,
				"DELETE":	DeleteRequest,
				"SET_STATUS": SetStatus,
				"GET_MY_PROPOSAL_FOR_A_REQUEST": GetMyProposalForARequest,
			},
			"rides": {
				"READ_ALL": GetAllRides,
				"READ_ONE": GetSingleRide,
				"CREATE":	CreateRide,
				"UPDATE":	UpdateRide,
				"DELETE":	DeleteRide,
				"ACTIVE_RIDE": ActiveRide,
			},
			"proposals": {
				"GET_ALL": GetAllProposals,
				"GET_ALL_MY": GetAllMyProposals,
			},
			"chats": {
				"CREATE": CreateChatRoom,
			},
			"messages": {
				"CREATE": SendMessage,
			},
		},
		"img->auth": {
			"files": {
				"VERIFY_TOKEN": VerifyToken,
				"GET_CLAIMS":   GetClaims,
			},
		},
	}

	return router[topic][pubsubMessage.Entity][pubsubMessage.Operation](pubsubMessage)
}
