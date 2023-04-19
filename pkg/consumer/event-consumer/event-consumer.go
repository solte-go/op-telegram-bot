package event_consumer

import (
	"context"
	"go.uber.org/zap"
	"telegram-bot/solte.lab/pkg/events"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
	logger    *zap.Logger
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int, logger *zap.Logger) *Consumer {
	return &Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
		logger:    logger,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			c.logger.Warn("received context done, stopping consumer")
			return ctx.Err()
		default:
		}

		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			c.logger.Error("can't fetch events", zap.Error(err))
			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		if err := c.eventsHandler(gotEvents); err != nil {
			c.logger.Warn("can't handle events", zap.Error(err))
			continue
		}
	}
}

//todo: 1. Data recovery, fallbacks (retry, save to memory, return to db).
//todo: 2 batch processing errors

func (c *Consumer) eventsHandler(events []events.Event) error {
	for _, event := range events {
		//c.logger.Info(fmt.Sprintf("Got new event: %s", event.Text))

		if err := c.processor.Process(event); err != nil {
			c.logger.Error("can't process event", zap.Error(err))
			continue
		}
	}
	return nil
}
