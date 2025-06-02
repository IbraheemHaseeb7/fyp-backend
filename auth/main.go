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
	utils.Requests = utils.NewSafeMap()

	// creating a publisher
	publisher, err := pubsub.NewPublisher(&pubsub.Publisher{
		URI:   amqpURI,
		Queue: "auth_db",
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	// subscribing to listen for responses from db server
	auth2dbSubscriber, err := pubsub.NewSubscriber(&pubsub.Subscriber{
		URI:   amqpURI,
		Queue: "auth_db",
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	auth2dbSubscriber.ConsumeMessages(
		"db->auth",
		func(pm pubsub.PubsubMessage) {
			utils.Requests.Load(pm.UUID) <- pm
		},
		handler.Handle,
	)

	// subscribing to listen for responses from db server
	img2authSubscriber, err := pubsub.NewSubscriber(&pubsub.Subscriber{
		URI:   amqpURI,
		Queue: "img_auth",
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	img2authSubscriber.ConsumeMessages(
		"img->auth",
		func(pm pubsub.PubsubMessage) {
			pm.Topic = "auth->img"
			publisher.PublishMessage(pm)
		},
		handler.Handle,
	)

	// starting HTTP server
	http.StartHTTPServer(publisher)
}
