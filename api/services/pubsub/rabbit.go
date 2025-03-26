package pubsub

import (

	// "github.com/rs/zerolog"
	// "retroleague.org/api/internal/env"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"
)

type RabbitMQService struct {
	publisher  *amqp.Publisher
	subscriber *amqp.Subscriber
}

func (rmq RabbitMQService) GetPublisher() *amqp.Publisher {
	return rmq.publisher
}

func (rmq RabbitMQService) GetSubscriber() *amqp.Subscriber {
	return rmq.subscriber
}

func (rmq RabbitMQService) Close() {
	rmq.publisher.Close()
	rmq.subscriber.Close()
}

func NewRabbitMQService(rabbitURI string) (*RabbitMQService, error) {
	cfg := amqp.NewDurableQueueConfig(rabbitURI)

	publisher, err := amqp.NewPublisher(cfg, watermill.NewStdLogger(false, false))
	if err != nil {
		return nil, err
	}

	subscriber, err := amqp.NewSubscriber(cfg, watermill.NewStdLogger(false, false))
	if err != nil {
		return nil, err
	}

	return &RabbitMQService{publisher, subscriber}, nil
}
