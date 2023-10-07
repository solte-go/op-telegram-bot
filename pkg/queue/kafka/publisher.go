package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"telegram-bot/solte.lab/pkg/config"
	"telegram-bot/solte.lab/pkg/models"

	"time"

	"go.uber.org/zap"
)

const (
	PlantMessage   = "plant"
	UserMessage    = "user"
	UnknownMessage = "unknown"
)

var connectionPool = make(map[string]*Publisher)

type Publisher struct {
	conf       *config.KafkaProducer
	producer   *kafka.Producer
	logger     *zap.Logger
	resultChan chan sendResult
}

type sendResult struct {
	err error
}

type messageWrapper struct {
	msg *kafka.Message
	idx int
}

func GetKafkaPublisher(connName string) (*Publisher, error) {
	if publisher, ok := connectionPool[connName]; ok {
		return publisher, nil
	}

	return nil, ErrNoConnectionWithProvidedName
}

func NewKafkaPublisher(ctx context.Context, conf *config.KafkaProducer) (*Publisher, error) {
	var (
		producer        *kafka.Producer
		err             error
		publisherLogger = zap.L().Named("kafka_publisher")
		delay           = 1 * time.Second
		attempt         = 0
	)

	if conf.ConnectionName == "" {
		return nil, ErrEmptyConnectionName
	}

	ConnectionCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for {
		select {
		case <-ConnectionCtx.Done():
			return nil, ConnectionCtx.Err()
		case <-time.After(delay):
			attempt++
			delay *= 2
			publisherLogger.Info("Try connect to v1", zap.Int("attempt", attempt))

			producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": conf.Brokers})
			if err == nil {
				p := &Publisher{
					resultChan: make(chan sendResult),
					conf:       conf,
					producer:   producer,
					logger:     publisherLogger,
				}
				connectionPool[conf.ConnectionName] = p
				go p.errorHandling(ctx)
				p.logger.Info("Kafka publisher connected", zap.String("connection_name", conf.ConnectionName))
				return p, err
			}
		}
	}
}

func (p *Publisher) Close() {
	p.producer.Close()
}

func (p *Publisher) errorHandling(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			p.logger.Info("Publisher error handling stopped")
			close(p.resultChan)
			return

		case result := <-p.resultChan:
			if result.err != nil {
				p.logger.Error("Error while sending message", zap.Error(result.err))
			}
		}
	}
}

func (p *Publisher) PrepareMessage(message interface{}) (*kafka.Message, error) {
	switch m := message.(type) {
	case *models.User:
		jsn, err := json.Marshal(m)
		if err != nil {
			return nil, err
		}

		fmt.Println(m)

		return &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &p.conf.Topic, Partition: kafka.PartitionAny},
			Key:            []byte(UserMessage),
			Value:          jsn,
		}, nil
	case string:
		return &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &p.conf.Topic, Partition: kafka.PartitionAny},
			Key:            []byte(PlantMessage),
			Value:          []byte(m),
		}, nil
	default:
		return nil, ErrUnknownMessageType
	}
}

func (p *Publisher) SendMessage(ctx context.Context, message *kafka.Message) {
	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)

	err := p.producer.Produce(message, deliveryChan)
	if err != nil {
		p.resultChan <- sendResult{err: err}
		return
	}

	select {
	case <-ctx.Done():
		p.resultChan <- sendResult{err: ctx.Err()}
		<-deliveryChan
		return

	case event := <-deliveryChan:
		msg := event.(*kafka.Message)
		if msg.TopicPartition.Error != nil {
			p.resultChan <- sendResult{err: msg.TopicPartition.Error}
			return
		}

		p.resultChan <- sendResult{err: nil}
	}
}

//func (p *Publisher) SendMessages(ctx context.Context, messages []*kafka.Message) []error {
//	resultChan := make(chan sendResult, len(messages))
//	defer close(resultChan)
//
//	for index, message := range messages {
//		go p.sendMessage(ctx, messageWrapper{msg: message, idx: index}, resultChan)
//		fmt.Printf("Send message: %s", message.Value)
//		time.Sleep(1 * time.Second)
//	}
//
//	result := make([]error, len(messages))
//	for range messages {
//		receiver := <-resultChan
//		result[receiver.idx] = receiver.err
//	}
//
//	return result
//}
