package main

import (
	"context"
	"flag"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	tgClient "telegram-bot/solte.lab/pkg/clients/telegram"
	"telegram-bot/solte.lab/pkg/config"
	eventConsumer "telegram-bot/solte.lab/pkg/consumer/eventconsumer"
	"telegram-bot/solte.lab/pkg/events/telegram"
	"telegram-bot/solte.lab/pkg/logging"
	"telegram-bot/solte.lab/pkg/queue"
)

var env string

const (
	batchSize   = 100
	storagePath = "data"
)

func init() {
	flag.StringVar(&env, "env", "dev", "Environment (dev, prod)")
	flag.Parse()
}

func main() {
	conf, err := config.LoadConf(env)
	if err != nil {
		panic(err)
	}

	logger, err := logging.NewLogger(conf.Logging)
	if err != nil {
		panic(err)
	}

	logger.Debug("Telegram Bot Started")
	defer logger.Sync() //nolint

	undo := zap.ReplaceGlobals(logger)
	defer undo()

	if conf.TG.Token == "" {
		logger.Fatal("Telegram token is empty")
	}

	tg := tgClient.New(conf.TG.Host, conf.TG.Token)

	ctx := waitQuitSignal(context.Background())

	k, err := queue.NewKafkaPublisher(ctx, conf.KafkaProducer)
	if err != nil {
		logger.Fatal("Kafka publisher error", zap.Error(err))
	}

	eventProcessor := telegram.New(tg, k, logger)

	c := eventConsumer.New(eventProcessor, eventProcessor, batchSize, logger)
	if err := c.Start(ctx); err != nil {
		logger.Warn("bot is shutting down", zap.String("reason", err.Error()))
	}
}

func waitQuitSignal(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		quit := make(chan os.Signal, 1)
		defer close(quit)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		<-quit
		cancel()
	}()

	return ctx
}
