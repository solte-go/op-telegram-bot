package kafka

//import (
//	"context"
//	"fmt"
//	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
//	"go.uber.org/zap"
//	"math"
//	"telegram-bot/solte.lab/pkg/config"
//	"time"
//)
//
//type RawMessage struct {
//	Key   []byte
//	Value []byte
//	//Traceparent []byte
//}
//
//type Consumer struct {
//	config     *config.KafkaConsumer
//	client     *kafka.Consumer
//	logger     *zap.Logger
//	partitions []kafka.TopicPartition
//	toCommit   map[int32]kafka.TopicPartition
//}
//
//func New(conf *config.KafkaConsumer) (*Consumer, error) {
//	var err error
//	var consumer *kafka.Consumer
//
//	logger := zap.L().Named("consumer")
//	for i := 0; ; i++ {
//		logger.Info("Try connect to kafka", zap.Int("attempt", i))
//
//		consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
//			"bootstrap.servers":               conf.Brokers,
//			"group.id":                        conf.Group,
//			"auto.offset.reset":               conf.AutoOffsetReset,
//			"enable.auto.commit":              false,
//			"enable.auto.offset.store":        false,
//			"session.timeout.ms":              conf.SessionTimeout,
//			"go.application.rebalance.enable": true,
//		})
//
//		if err == nil {
//			logger.Info("Connected")
//			break
//		}
//
//		logger.Error("Cant connect", zap.Error(err))
//		if i > conf.ConnectionRetries {
//			break
//		}
//
//		time.Sleep(time.Second)
//	}
//
//	if err != nil {
//		return nil, err
//	}
//
//	c := &Consumer{
//		config:   conf,
//		client:   consumer,
//		logger:   logger,
//		toCommit: map[int32]kafka.TopicPartition{},
//	}
//
//	err = c.client.SubscribeTopics([]string{conf.Topic}, c.rebalanceCB)
//
//	return c, err
//}
//
//func (c *Consumer) rebalanceCB(kc *kafka.Consumer, ev kafka.Event) error {
//	c.logger.Info("RebalanceCb executed")
//	switch e := ev.(type) {
//	case kafka.Error:
//		c.logger.Error("Rebalance kafka.Error", zap.Error(e))
//		return e
//	case kafka.AssignedPartitions:
//		c.logger.Info("Rebalance", zap.Any("event", e))
//		c.partitions = e.Partitions
//		err := c.CommitOffsets()
//		if err != nil {
//			c.logger.Error("Commit error", zap.Error(err))
//		}
//
//		err = c.client.Assign(e.Partitions)
//		if err != nil {
//			c.logger.Error("Assign() error", zap.Error(err))
//		}
//	case kafka.RevokedPartitions:
//		c.logger.Info("Rebalance", zap.Any("event", e))
//		err := c.CommitOffsets()
//		if err != nil {
//			c.logger.Error("Commit error", zap.Error(err))
//		}
//
//		err = c.client.Unassign()
//		if err != nil {
//			c.logger.Error("Unassign() error", zap.Error(err))
//		}
//	default:
//		// Ignore other event types
//		c.logger.Debug("Rebalance Ignored", zap.Any("event", e))
//	}
//
//	return nil
//}
//
//func (c *Consumer) CommitOffsets() error {
//	if len(c.toCommit) == 0 {
//		c.logger.Debug("Nothing to commit")
//		return nil
//	}
//
//	offsets := make([]kafka.TopicPartition, 0, len(c.toCommit))
//	for _, tp := range c.toCommit {
//		tp.Offset++
//		offsets = append(offsets, tp)
//	}
//
//	var err error
//	var commited []kafka.TopicPartition
//	for try := 0; try <= c.config.OffsetCommitRetries; try++ {
//		c.logger.Debug("About to commit offsets", zap.Int("try", try))
//		commited, err = c.client.CommitOffsets(offsets)
//		if err == nil {
//			break
//		}
//	}
//
//	if err != nil {
//		c.logger.Error(
//			"Can't commit offsets",
//			zap.Int("retries", c.config.OffsetCommitRetries),
//			zap.Error(err),
//		)
//
//		return err
//	}
//
//	c.toCommit = make(map[int32]kafka.TopicPartition)
//
//	for _, offset := range commited {
//		c.logger.Debug("Offsets committed", zap.String("offset", offset.String()))
//	}
//
//	return nil
//}
//
//func (c *Consumer) readMessage() (*kafka.Message, error) {
//	timeout := c.config.PollTimeout
//
//	var absTimeout time.Time
//	var timeoutMs int
//
//	if timeout > 0 {
//		absTimeout = time.Now().Add(timeout)
//		timeoutMs = (int)(timeout.Seconds() * 1000.0)
//	} else {
//		timeoutMs = (int)(timeout)
//	}
//
//	for {
//		ev := c.client.Poll(timeoutMs)
//
//		switch e := ev.(type) {
//		case *kafka.Message:
//			if e.TopicPartition.Error != nil {
//				return e, e.TopicPartition.Error
//			}
//			return e, nil
//		case kafka.Error:
//			return nil, e
//		case kafka.PartitionEOF:
//			c.logger.Info("PartitionEOF")
//		default:
//			// Ignore other event types
//			if e != nil {
//				c.logger.Debug(
//					fmt.Sprintf("Unexpected message type: %v", e),
//				)
//			}
//		}
//
//		if timeout > 0 {
//			// Calculate remaining time
//			timeoutMs = int(math.Max(0.0, absTimeout.Sub(time.Now()).Seconds()*1000.0))
//		}
//
//		if timeoutMs == 0 && ev == nil {
//			c.logger.Debug("ReadMessage timed out")
//			return nil, nil
//		}
//	}
//}
//
//func (c *Consumer) ReadBatch() []RawMessage {
//	batch := make([]RawMessage, 0, c.config.BatchSize)
//
//	currentBatchSize := 0
//	timer := time.NewTimer(c.config.PollTimeout)
//BatchLoop:
//	for currentBatchSize < c.config.BatchSize {
//		select {
//		case <-timer.C:
//			c.logger.Info("Batch reading timeout")
//			break BatchLoop
//		default:
//			c.logger.Debug("Reading next batch message")
//			m, err := c.readMessage()
//			if m != nil {
//				c.toCommit[m.TopicPartition.Partition] = m.TopicPartition
//				currentBatchSize++
//			} else {
//				c.logger.Debug("No new messages")
//				continue
//			}
//			if err != nil {
//				c.logger.Error("Can't read event", zap.Error(err))
//				continue
//			}
//
//			//var traceparent []byte
//			//for _, h := range m.Headers {
//			//	if h.Key == "elasticapmtraceparent" || h.Key == "traceparent" {
//			//		traceparent = h.Value
//			//		break
//			//	}
//			//}
//
//			c.logger.Debug("Processed offset " + m.TopicPartition.String())
//			batch = append(batch, RawMessage{Key: m.Key, Value: m.Value})
//		}
//	}
//
//	return batch
//}
//
//func (c *Consumer) PollMessages(ctx context.Context, timeoutMs int, messageChan chan *kafka.Message) {
//	commitMessages := 0
//
//	for {
//		select {
//		case <-ctx.Done():
//			c.logger.Info("Context done", zap.String("Reason", ctx.Err().Error()))
//			_, err := c.client.Commit()
//			if err != nil {
//				c.logger.Error("Error commit client", zap.Error(err))
//			}
//			err = c.client.Close()
//			if err != nil {
//				c.logger.Error("Error closing client", zap.Error(err))
//				return
//			}
//			return
//
//		default:
//		}
//
//		ev := c.client.Poll(timeoutMs)
//		if ev == nil {
//			continue
//		}
//
//		switch e := ev.(type) {
//		case *kafka.Message:
//			messageChan <- e
//			commitMessages++
//			if commitMessages == 10 {
//				go func() {
//					c.client.Commit()
//				}()
//				//_, err := c.client.CommitMessage(e)
//				//if err != nil {
//				//	c.logger.Error("Error committing message", zap.Error(err))
//				//	return
//				//}
//			}
//
//			_, err := c.client.StoreMessage(e)
//			if err != nil {
//				c.logger.Error(
//					"Error storing offset after message",
//					zap.Error(err),
//					zap.String("TopicPartition", e.TopicPartition.String()),
//				)
//			}
//		case kafka.Error:
//			c.logger.Error("Error", zap.Error(e), zap.String("Code", e.Code().String()))
//			if e.Code() == kafka.ErrAllBrokersDown {
//				ctx.Done()
//				return
//			}
//		default:
//			c.logger.Info("Ignored", zap.Any("Event", e))
//		}
//	}
//}
