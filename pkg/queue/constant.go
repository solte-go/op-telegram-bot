package queue

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"os"
)

func (c *Consumer) PollMessages(ctx context.Context, timeoutMs int, messageChan chan *kafka.Message) {
consume:
	for {
		select {
		case sig := <-ctx.Done():
			fmt.Printf("Caught signal %v: terminating\n", sig)
			break consume

		default:
			ev := c.client.Poll(timeoutMs)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("%% Message on %s:\n%s\n",
					e.TopicPartition, string(e.Value))
				if e.Headers != nil {
					fmt.Printf("%% Headers: %v\n", e.Headers)
				}
				messageChan <- e

				_, err := c.client.StoreMessage(e)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%% Error storing offset after message %s:\n",
						e.TopicPartition)
					fmt.Println(err)
				}
			case kafka.Error:
				// Errors should generally be considered
				// informational, the client will try to
				// automatically recover.
				// But in this example we choose to terminate
				// the application if all brokers are down.
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {

				}
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}
}
