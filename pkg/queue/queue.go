package queue

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Consumer interface {
	PollMessages(ctx context.Context, timeoutMs int, messageChan chan *kafka.Message)
}

type Producer interface {
	PrepareMessage(message interface{}) (*kafka.Message, error)
	SendMessage(ctx context.Context, message *kafka.Message)
}
