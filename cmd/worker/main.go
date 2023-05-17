package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"telegram-bot/solte.lab/pkg/api"
	"telegram-bot/solte.lab/pkg/api/handlers/metrics"
	tgClient "telegram-bot/solte.lab/pkg/clients/telegram"
	"telegram-bot/solte.lab/pkg/config"
	"telegram-bot/solte.lab/pkg/events/telegram"
	"telegram-bot/solte.lab/pkg/logging"
	"telegram-bot/solte.lab/pkg/queue/kafka"
	"telegram-bot/solte.lab/pkg/storage/emsql"
	"telegram-bot/solte.lab/pkg/storage/postgresql"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
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

	storage, err := postgresql.New(conf.Postgres)
	if err != nil {
		logger.Fatal("can't initialize storage", zap.Error(err))
	}

	consumer, err := kafka.New(conf.KafkaConsumer)
	if err != nil {
		logger.Fatal("can't initialize consumer", zap.Error(err))
	}

	server := api.New(logger)
	go server.Run(ctx, conf.APIs.Worker.Port, &metrics.Worker{})

	responder := telegram.NewResponder(ctx, tg, emsql.New(storage), consumer, logger)

	responder.Run(ctx)
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
