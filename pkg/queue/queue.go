package queue

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Consumer interface {
	Run(ctx context.Context, ch chan *kafka.Message) error
}

type Producer interface {
	PrepareMessage(message interface{}) (*kafka.Message, error)
	SendMessage(ctx context.Context, message *kafka.Message)
}
