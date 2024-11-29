package broker

import (
	"encoding/json"
	"fmt"
	"micro/handler"
	"micro/utils"

	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

func RabbitMQ(db *gorm.DB) {
	// connect to rabbit mq
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	utils.ErrorHandler(err)
	defer conn.Close()

	// creating channel
	ch, err := conn.Channel()
	utils.ErrorHandler(err)
	defer ch.Close()

	// creating queues
	reqQ, err := ch.QueueDeclare("data_request_queue", false, false, false, false, nil)
	utils.ErrorHandler(err)
	resQ, err := ch.QueueDeclare("data_response_queue", false, false, false, false, nil)
	utils.ErrorHandler(err)

	// setting up consumer
	msgs, err := ch.Consume(reqQ.Name, "", false, false, false, false, nil)
	utils.ErrorHandler(err)

	fmt.Println("Waitin for messages...")

	forever := make(chan bool)
	for msg := range msgs {

		var req handler.Request
		err := json.Unmarshal([]byte(msg.Body), &req)
		utils.ErrorHandler(err)

		// sending response back to go service	
		response := handler.GlobalHandler(req, db).(string)
		err = ch.Publish("", resQ.Name, false, false, amqp.Publishing {
			ContentType: "text/plain",
			Body: []byte(response),
			CorrelationId: msg.CorrelationId,
		}) 

		msg.Ack(false)
		utils.ErrorHandler(err)
	}
	<-forever
}

