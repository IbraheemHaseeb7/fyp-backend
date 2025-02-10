package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
)

type PubsubMessage struct {
	Operation string `json:"operation"`
	Payload   any    `json:"payload"`
	Entity    string `json:"entity"`
	Instance  *message.Message
	Topic     string `json:"topic"`
	UUID	  string `json:"uuid"`
}

type Publisher struct {
	URI      string
	Queue    string
	Instance *amqp.Publisher
}

type Subscriber struct {
	URI      string
	Queue    string
	Instance *amqp.Subscriber
}

var Service string

func (pm *PubsubMessage) NewPubsubMessage(p *amqp.Publisher) (*PubsubMessage, error) {
	if pm.UUID == "" {
		pm.UUID = watermill.NewUUID()
	}

	payload, err := json.Marshal(pm)
	if err != nil {
		return nil, err
	}

	pm.Instance = message.NewMessage(pm.UUID, payload)

	return pm, nil
}

func NewPublisher(p *Publisher) (*Publisher, error) {
	if p.URI == "" {
		p.URI = "amqp://guest:guest@localhost:5672/"
	}

	// RabbitMQ config
	amqpConfig := amqp.NewDurableQueueConfig(p.URI)

	// Create publisher
	publisher, err := amqp.NewPublisher(amqpConfig, watermill.NewStdLogger(false, false))
	if err != nil {
		return nil, err
	}

	p.Instance = publisher

	return p, nil
}

func NewSubscriber(s *Subscriber) (*Subscriber, error) {

	amqpConfig := amqp.NewDurableQueueConfig(s.URI)

	subscriber, err := amqp.NewSubscriber(
		amqpConfig,
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		return nil, err
	}

	s.Instance = subscriber

	return s, nil
}

func (p *Publisher) PublishMessage(pm PubsubMessage) error {

	message, err := pm.NewPubsubMessage(p.Instance)
	if err != nil {
		return err
	}

	err = p.Instance.Publish(pm.Topic, message.Instance)
	if err != nil {
		return err
	}

	return nil
}

func (s *Subscriber) ConsumeMessages(topic string, function func(PubsubMessage), handle func(string, PubsubMessage) (PubsubMessage, error)) (any, error) {
	messages, err := s.Instance.Subscribe(context.Background(), topic)
	if err != nil {
		return nil, err
	}

	go func() {
		for message := range messages {

			var pubsubMessage PubsubMessage
			if err := json.Unmarshal(message.Payload, &pubsubMessage); err != nil {
				fmt.Println(err.Error())
				message.Ack()
				continue
			}

			recMsg, err := handle(topic, pubsubMessage)
			if err != nil {
				fmt.Println(err.Error())
				message.Ack()
				continue
			}

			if function != nil {
				function(recMsg)
			}
			message.Ack()
		}
	}()
	return nil, nil
}

