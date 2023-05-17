package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"
	"telegram-bot/solte.lab/pkg/clients/telegram"
	"telegram-bot/solte.lab/pkg/models"
	"telegram-bot/solte.lab/pkg/queue"
	v1 "telegram-bot/solte.lab/pkg/queue/kafka"
	"telegram-bot/solte.lab/pkg/storage/emsql"
	"telegram-bot/solte.lab/pkg/worker"
)

type Responder struct {
	tg       *telegram.Client
	consumer queue.Consumer
	//storage  storage.Storage
	logger *zap.Logger
	worker *worker.Worker
}

func NewResponder(
	ctx context.Context,
	client *telegram.Client,
	storage emsql.OPContract,
	kafka *v1.Consumer,
	logger *zap.Logger,
) *Responder {
	return &Responder{
		tg:       client,
		consumer: kafka,
		worker:   worker.New(ctx, storage),
		logger:   logger,
	}
}

func (r *Responder) Run(ctx context.Context) {
	ch := make(chan *kafka.Message)

	go r.consumer.PollMessages(ctx, 100, ch)

	for {
		select {
		case <-ctx.Done():
			r.logger.Info("telegram responder stopped")
			return
		case msg := <-ch:
			switch string(msg.Key) {
			case v1.PlantMessage:
				fmt.Println("plant message: ", string(msg.Value))
			case v1.UserMessage:
				user, err := r.Decode(msg.Value)
				if err != nil {
					r.logger.Error("can't decode user message", zap.Error(err))
					continue
				}
				err = r.doCmd(&user)
				if err != nil {
					r.logger.Error("can't process user command", zap.Error(err))
					continue
				}
			}
		}
	}
}

func (r *Responder) Decode(data []byte) (models.User, error) {
	var user models.User
	err := json.Unmarshal(data, &user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
