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
	utils.Requests = utils.NewSafeMap()

	publisher, err := pubsub.NewPublisher(&pubsub.Publisher{
		URI:   os.Getenv("AMQP_STRING"),
		Queue: "img_auth",
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	// subscribing to listen for responses from db server
	img2authSubscriber, err := pubsub.NewSubscriber(&pubsub.Subscriber{
		URI:   os.Getenv("AMQP_STRING"),
		Queue: "img_auth",
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	img2authSubscriber.ConsumeMessages(
		"auth->img",
		func(pm pubsub.PubsubMessage) {
			utils.Requests.Load(pm.UUID) <- pm
		},
		handler.Handle,
	)

	img2dbSubscriber, err := pubsub.NewSubscriber(&pubsub.Subscriber{
		URI:   os.Getenv("AMQP_STRING"),
		Queue: "img_db",
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	img2dbSubscriber.ConsumeMessages(
		"db->img",
		func(pm pubsub.PubsubMessage) {
			utils.Requests.Load(pm.UUID) <- pm
		},
		handler.Handle,
	)

	http.StartHTTPServer(publisher)
}
