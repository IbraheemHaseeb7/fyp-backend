package main

import (
	"fmt"
	"os"

	"github.com/IbraheemHaseeb7/fyp-backend/handler"
	"github.com/IbraheemHaseeb7/fyp-backend/http"
	"github.com/IbraheemHaseeb7/fyp-backend/utils"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()
	amqpURI := os.Getenv("AMQP_STRING")

	// creating a publisher
	publisher, err := pubsub.NewPublisher(&pubsub.Publisher{
		URI:   amqpURI,
		Queue: "auth_db",
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	// subscribing to listen for responses from db server
	subscriber, err := pubsub.NewSubscriber(&pubsub.Subscriber{
		URI:   amqpURI,
		Queue: "auth_db",
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	subscriber.ConsumeMessages(
		"db->auth",
		func(pm pubsub.PubsubMessage) {
			utils.Requests[pm.UUID] <- pm
		},
		handler.Handle,
	)

	// starting HTTP server
	http.StartHTTPServer(publisher)
}
