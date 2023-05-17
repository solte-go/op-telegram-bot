package main

import (
	"context"
	"flag"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"telegram-bot/solte.lab/pkg/api"
	"telegram-bot/solte.lab/pkg/api/handlers/admin"
	"telegram-bot/solte.lab/pkg/api/handlers/authorized"
	"telegram-bot/solte.lab/pkg/config"
	"telegram-bot/solte.lab/pkg/logging"
	"telegram-bot/solte.lab/pkg/storage/emsql"
	"telegram-bot/solte.lab/pkg/storage/postgresql"
)

var env string

func init() {
	flag.StringVar(&env, "env", "dev", "Environment (dev, prod)")
	flag.Parse()
}

// api server template
func main() {
	conf, err := config.LoadConf(env)
	if err != nil {
		panic(err)
	}

	logger, err := logging.NewLogger(conf.Logging)
	if err != nil {
		panic(err)
	}

	logger.Debug("API Server Started")
	defer logger.Sync() //nolint

	undo := zap.ReplaceGlobals(logger)
	defer undo()

	ctx := waitQuitSignal(context.Background())

	db, err := postgresql.New(conf.Postgres)
	if err != nil {
		logger.Fatal("failed to init storage", zap.Error(err))
	}

	server := api.New(logger)
	server.Run(
		ctx, conf.APIs.UI.Port,
		admin.New(emsql.New(db)),
		authorized.New(emsql.New(db), conf.APIs.UI.StaticContentPath),
	)
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
