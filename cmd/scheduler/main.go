package main

import (
	"context"
	"flag"

	"go.uber.org/zap"
	"telegram-bot/solte.lab/pkg/config"
	"telegram-bot/solte.lab/pkg/logging"
	"telegram-bot/solte.lab/pkg/queue/kafka"
	"telegram-bot/solte.lab/pkg/scheduler"
	"telegram-bot/solte.lab/pkg/storage/postgresql"

	_ "github.com/lib/pq"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	db, err := postgresql.New(conf.Postgres)
	if err != nil {
		logger.Fatal("failed to init storage", zap.Error(err))
	}

	k, err := kafka.NewKafkaPublisher(ctx, conf.KafkaProducer)
	if err != nil {
		logger.Fatal("Kafka publisher error", zap.Error(err))
	}

	sch := scheduler.New(db, k, logger, conf.Scheduler)

	err = sch.Run(ctx)
	if err != nil {
		return
	}
}
