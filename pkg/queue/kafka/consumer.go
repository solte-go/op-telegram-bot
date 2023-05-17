package kafka

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"
	"telegram-bot/solte.lab/pkg/config"
	"time"
)

type Consumer struct {
	config     *config.KafkaConsumer
	client     *kafka.Consumer
	logger     *zap.Logger
	partitions []kafka.TopicPartition
	channel    chan *kafka.Message
}

func New(conf *config.KafkaConsumer) (*Consumer, error) {
	var err error
	var consumer *kafka.Consumer

	logger := zap.L().Named("consumer")
	for i := 0; ; i++ {
		logger.Info("Try connect to kafka", zap.Int("attempt", i))

		consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers":               conf.Brokers,
			"group.id":                        conf.Group,
			"auto.offset.reset":               conf.AutoOffsetReset,
			"enable.auto.commit":              false,
			"session.timeout.ms":              conf.SessionTimeout,
			"go.application.rebalance.enable": true,
			//"enable.auto.offset.store":        false,
		})

		if err == nil {
			logger.Info("Connected")
			break
		}

		logger.Error("Cant connect", zap.Error(err))
		if i > conf.ConnectionRetries {
			break
		}

		time.Sleep(time.Second)
	}

	if err != nil {
		return nil, err
	}

	c := &Consumer{
		config: conf,
		client: consumer,
		logger: logger,
	}

	err = c.client.SubscribeTopics([]string{conf.Topic}, c.rebalanceCallback)
	return c, err
}

func (c *Consumer) Run(ctx context.Context, ch chan *kafka.Message) error {

	c.channel = ch

	run := true
	for run == true {
		select {
		case <-ctx.Done():
			c.logger.Info("Revise context done: Closing consumer")
			run = false
		default:
			ev := c.client.Poll(100)
			if ev == nil {
				continue
			}

			if err := c.processEvent(ev); err != nil {
				c.logger.Error("Failed to process event", zap.Error(err))
			}
		}
	}

	c.logger.Info("Closing consumer")
	err := c.client.Close()
	if err != nil {
		c.logger.Error("Failed to close consumer", zap.Error(err))
		return err
	}

	return nil
}

// processEvent processes the message/error received from the kafka Consumer's
// Poll() method.
func (c *Consumer) processEvent(ev kafka.Event) error {
	switch e := ev.(type) {

	case *kafka.Message:
		c.channel <- e

		if err := c.maybeCommit(e.TopicPartition); err != nil {
			return err
		}

	case kafka.Error:
		// Errors should generally be considered informational, the client
		// will try to automatically recover.
		c.logger.Error("Kafka error", zap.Error(e))

	default:
		fmt.Printf("Ignored %v\n", e)
	}

	return nil
}

// maybeCommit is called for each message we receive from a Kafka topic.
// This method can be used to apply some arbitary logic/processing to the
// offsets, write the offsets into some external storage, and finally, to
// decide when we want to commit already-stored offsets into Kafka.
func (c *Consumer) maybeCommit(topicPartition kafka.TopicPartition) error {
	//additional logic to commit offset should be here if needed
	//if topicPartition.Offset%10 != 0 {
	//	return nil
	//}

	commitedOffsets, err := c.client.Commit()

	// ErrNoOffset occurs when there are no stored offsets to commit. This
	// can happen if we haven't stored anything since the last commit.
	// While this will never happen for this example since we call this method
	// per-message, and thus, always have something to commit, the error
	// handling is illustrative of how to handle it in cases we call Commit()
	// in another way, for example, every N seconds.
	if err != nil && err.(kafka.Error).Code() != kafka.ErrNoOffset {
		return err
	}

	c.logger.Debug("Committing offsets", zap.Any("offsets", commitedOffsets))
	return nil
}

// rebalanceCallback is called on each group rebalance to assign additional
// partitions, or remove existing partitions, from the consumer's current
// assignment.
//
// A rebalance occurs when a consumer joins or leaves a consumer group, if it
// changes the topic(s) it's subscribed to, or if there's a change in one of
// the topics it's subscribed to, for example, the total number of partitions
// increases.
//
// The application may use this optional callback to inspect the assignment,
// alter the initial start offset (the .Offset field of each assigned partition),
// and read/write offsets to commit to an alternative store outside of Kafka.
func (c *Consumer) rebalanceCallback(kc *kafka.Consumer, event kafka.Event) error {
	switch ev := event.(type) {
	case kafka.AssignedPartitions:
		fmt.Printf("%% %s rebalance: %d new partition(s) assigned: %v\n",
			c.client.GetRebalanceProtocol(), len(ev.Partitions), ev.Partitions)

		// Assign the partitions to the consumer and setup offsets here if needed.

		err := c.client.Assign(ev.Partitions)
		if err != nil {
			return err
		}

	case kafka.RevokedPartitions:
		fmt.Printf("%% %s rebalance: %d partition(s) revoked: %v\n",
			c.client.GetRebalanceProtocol(), len(ev.Partitions), ev.Partitions)

		// Usually, the rebalance callback for `RevokedPartitions` is called
		// just before the partitions are revoked. We can be certain that a
		// partition being revoked is not yet owned by any other consumer.
		// This way, logic like storing any pending offsets or committing
		// offsets can be handled.
		// However, there can be cases where the assignment is lost
		// involuntarily. In this case, the partition might already be owned
		// by another consumer, and operations including committing
		// offsets may not work.
		if c.client.AssignmentLost() {
			// Our consumer has been kicked out of the group and the
			// entire assignment is thus lost.
			c.logger.Error("Assignment lost involuntarily, commit may fail")
		}

		// Since enable.auto.commit is unset, we need to commit offsets manually
		// before the partition is revoked.
		committedOffsets, err := c.client.Commit()

		if err != nil && err.(kafka.Error).Code() != kafka.ErrNoOffset {
			c.logger.Error("Failed to commit offsets", zap.Error(err))
			return err
		}

		c.logger.Debug("Committing offsets", zap.Any("offsets", committedOffsets))

		// Similar to Assign, client automatically calls Unassign() unless the
		// callback has already called that method. Here, we don't call it.

	default:
		c.logger.Error("Unexpected event type", zap.Any("event", event))
	}

	return nil
}
