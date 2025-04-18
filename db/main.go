package main

import (
	"fmt"
	"os"

	"github.com/IbraheemHaseeb7/fyp-backend/db"
	"github.com/IbraheemHaseeb7/fyp-backend/handler"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()
	db.Connect()

	amqpURI := os.Getenv("AMQP_STRING")
	pubsub.Service = "db"

	subscriber, err := pubsub.NewSubscriber(&pubsub.Subscriber{
		URI:   amqpURI,
		Queue: "auth_db",
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	publisher, err := pubsub.NewPublisher(&pubsub.Publisher{
		URI:   amqpURI,
		Queue: "auth_db",
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	imgSubscriber, err := pubsub.NewSubscriber(&pubsub.Subscriber{
		URI:   amqpURI,
		Queue: "img_db",
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	imgPublisher, err := pubsub.NewPublisher(&pubsub.Publisher{
		URI:   amqpURI,
		Queue: "img_db",
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	subscriber.ConsumeMessages(
		"auth->db",
		func(pm pubsub.PubsubMessage) {
			pm.Topic = "db->auth"
			publisher.PublishMessage(pm)
		},
		handler.Handle,
	)

	imgSubscriber.ConsumeMessages(
		"img->db",
		func(pm pubsub.PubsubMessage) {
			pm.Topic = "db->img"
			imgPublisher.PublishMessage(pm)
		},
		handler.Handle,
	)
	select {}
}
