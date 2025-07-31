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
				"STORE_DEVICE_TOKEN": StoreDeviceToken,
				"RESET_PASSWORD": ResetPassword,
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
				"GET_MY_PROPOSAL_FOR_A_REQUEST": GetMyProposalForARequest,
				"GET_MATCHED_PROPOSAL_OF_A_REQUEST": GetMatchedProposalOfARequest,
				"GET_MATCHES": GetMatches,
				"GET_ACTIVE_REQUEST": GetActiveRequest,
			},
			"rides": {
				"READ_ALL": GetAllRides,
				"READ_ONE": GetSingleRide,
				"CREATE":	CreateRide,
				"UPDATE":	UpdateRide,
				"DELETE":	DeleteRide,
				"ACTIVE_RIDE":	ActiveRide,
			},
			"proposals": {
				"GET_ALL":	GetAllProposals,
				"GET_ALL_MY":	GetAllMyProposals,
			},
			"chats": {
				"CREATE":	CreateChatRoom,
				"GET_MESSAGES": GetChatMessages,
			},
			"messages": {
				"CREATE": SendMessage,
			},
			"rooms": {
				"READ_ONE": GetRoom,
				"READ_WITH_IDS": GetRoomWithIDs,
			},
			"feedback": {
				"CREATE": AddFeedback,
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
