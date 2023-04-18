package main

import (
	"flag"
	"go.uber.org/zap"
	"log"
	tgClient "telegram-bot/solte.lab/pkg/clients/telegram"
	"telegram-bot/solte.lab/pkg/config"
	eventConsumer "telegram-bot/solte.lab/pkg/consumer/event-consumer"
	"telegram-bot/solte.lab/pkg/events/telegram"
	"telegram-bot/solte.lab/pkg/logging"
	"telegram-bot/solte.lab/pkg/storage/cache"
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

	if conf.TG.Token == "" {
		logger.Fatal("Telegram token is empty")
	}

	tg := tgClient.New(conf.TG.Host, conf.TG.Token)

	s, err := cache.New(conf.PostgreSQL)
	if err != nil {
		logger.Fatal("can't initialize storage", zap.Error(err))
	}

	eventProcessor := telegram.New(tg, s, logger)

	logger.Info("Bot started")

	c := eventConsumer.New(eventProcessor, eventProcessor, batchSize, logger)
	if err := c.Start(); err != nil {
		log.Fatalf("can't start consumer: %v", err)
	}
}
