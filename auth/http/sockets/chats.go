package sockets

import (
	"fmt"
	"log"

	socketio "github.com/googollee/go-socket.io"
)

// setup a websocket connection to handle chat messages
func SetupSocket() *socketio.Server {

	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "join_private", func(s socketio.Conn, data map[string]string) {
		userID := data["userId"]
		peerID := data["peerId"]
		roomID := fmt.Sprintf("room-%s-%s", userID, peerID)

		fmt.Printf("🔗 User %s joining private room %s\n", userID, roomID)

		s.Join(roomID)
	})

	server.OnEvent("/", "private_message", func(s socketio.Conn, data map[string]string) {
		roomID := data["room"]
		message := data["message"]
		sender := data["sender"]

		fmt.Printf("📩 [%s] %s\n", roomID, message)

		// Send to all members in the room (including sender, or exclude if needed)
		server.BroadcastToRoom("/", roomID, "private_message", map[string]string{
			"sender":  sender,
			"message": message,
		})
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
