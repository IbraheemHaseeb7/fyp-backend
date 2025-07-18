package sockets

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IbraheemHaseeb7/fyp-backend/utils"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/ThreeDotsLabs/watermill"
	socketio "github.com/googollee/go-socket.io"
)

// setup a websocket connection to handle chat messages
func SetupSocket(p *pubsub.Publisher) *socketio.Server {

	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		return nil
	})

	server.OnEvent("/", "join_private", func(s socketio.Conn, data map[string]string) {
		requestID := data["request_id"]
		proposalID := data["proposal_id"]

		response, err := CreateChatRoom(map[string]any{
			"request_id":  requestID,
			"proposal_id": proposalID,
		}, p)

		if err != nil {
			s.Emit("error", map[string]string{
				"error": err.Error(),
			})
			return
		}
		if response["error"] != nil {
			return
		}
		roomID := response["data"].(map[string]any)["id"]
			
		s.Emit("joined_private", map[string]string{
			"room": fmt.Sprintf("%v", roomID),
		})
		s.Join(fmt.Sprintf("%v", roomID))
	})

	server.OnEvent("/", "private_message", func(s socketio.Conn, data map[string]string) {
		roomID := data["room"]
		message := data["message"]
		sender := data["sender"]

		if roomID == "" || message == "" || sender == "" {
			s.Emit("error", map[string]string{
				"error": "room, message and sender are required",
			})
			return
		}

		// Store in the database
		response, err := SendMessage(map[string]any{
			"room_id": roomID,
			"sender":  sender,
			"message": message,
		}, p)
		if err != nil {
			s.Emit("error", map[string]string{
				"error": err.Error(),
			})
			return
		}
		if response["error"] != nil {
			s.Emit("error", map[string]string{
				"error": response["error"].(string),
			})
			return
		}

		// Send to all members in the room (including sender, or exclude if needed)
		server.BroadcastToRoom("/", roomID, "private_message", map[string]string{
			"sender":  sender,
			"message": message,
		})
	})

	server.OnEvent("/", "user_locations", func(s socketio.Conn, data map[string]any) {
		roomID := data["room"].(string)
		userID := data["user"]

		if roomID == "" || userID == "" {
			s.Emit("error", map[string]string{
				"error": "room, user and locations are required",
			})
			return
		}
		// Broadcast user locations to all members in the room
		server.BroadcastToRoom("/", roomID, "user_locations", data)
	})

	server.OnEvent("/", "ride_status", func(s socketio.Conn, data map[string]any) {
		roomID := data["room"].(string)

		if roomID == "" {
			s.Emit("error", map[string]string{
				"error": "room is required",
			})
			return
		}

		payload, err := json.Marshal(data); if err != nil {
			s.Emit("error", map[string]string{
				"error": "Could not marshal data",
			})
			return
		}

		uuid := watermill.NewUUID()
		pubsubMessage := pubsub.PubsubMessage{
			Entity:    "rides",
			Operation: "READ_STATUS",
			Topic:     "auth->db",
			UUID:      uuid,
			Payload:   string(payload),
		}

		err = p.PublishMessage(pubsubMessage)
		if err != nil {
			s.Emit("error", map[string]string{
				"error": "Could not receive data from the database",
			})
			return
		}

		response := (<-utils.Requests.Load(pubsubMessage.UUID)).Payload.(map[string]any)
		utils.Requests.Delete(pubsubMessage.UUID)

		// Broadcast user locations to all members in the room
		server.BroadcastToRoom("/", roomID, "ride_status", response)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("disconnected:", s.ID(), "reason:", reason)
		s.Close()
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("Socket server error: %v", err)
		}
	}()

	return server
}

func CreateChatRoom(reqBody any, p *pubsub.Publisher) (map[string]any, error) {
	uuid := watermill.NewUUID()
	utils.Requests.Store(uuid, make(chan pubsub.PubsubMessage))

	payload, err := json.Marshal(reqBody); if err != nil {
		return nil, err
	}

	// publishing a read message
	pubsubMessage := pubsub.PubsubMessage{
		Entity:    "chats",
		Operation: "CREATE",
		Topic:     "auth->db",
		UUID:      uuid,
		Payload:   string(payload),
	}

	err = p.PublishMessage(pubsubMessage)
	if err != nil {
		return nil, err
	}

	response := (<-utils.Requests.Load(pubsubMessage.UUID)).Payload.(map[string]any)
	utils.Requests.Delete(pubsubMessage.UUID)

	return response, nil
}

func SendMessage(reqBody any, p *pubsub.Publisher) (map[string]any, error) {
	uuid := watermill.NewUUID()
	utils.Requests.Store(uuid, make(chan pubsub.PubsubMessage))

	payload, err := json.Marshal(reqBody); if err != nil {
		return nil, err
	}

	// publishing a read message
	pubsubMessage := pubsub.PubsubMessage{
		Entity:    "messages",
		Operation: "CREATE",
		Topic:     "auth->db",
		UUID:      uuid,
		Payload:   string(payload),
	}

	err = p.PublishMessage(pubsubMessage)
	if err != nil {
		return nil, err
	}

	response := (<-utils.Requests.Load(pubsubMessage.UUID)).Payload.(map[string]any)
	utils.Requests.Delete(pubsubMessage.UUID)

	return response, nil
}
