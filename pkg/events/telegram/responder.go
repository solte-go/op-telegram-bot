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
	"telegram-bot/solte.lab/pkg/storage"
	"telegram-bot/solte.lab/pkg/storage/storagewrapper"
)

type Responder struct {
	tg       *telegram.Client
	consumer *queue.Consumer
	storage  storage.Storage
	logger   *zap.Logger
}

func NewResponder(ctx context.Context, client *telegram.Client, s *storagewrapper.StorageCache, kafka *queue.Consumer, logger *zap.Logger) *Responder {
	return &Responder{
		tg:       client,
		consumer: kafka,
		storage:  s,
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
			case queue.PlantMessage:
				fmt.Println("plant message: ", string(msg.Value))
			case queue.UserMessage:
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
	fmt.Println("user: ", user.Cmd)
	return user, nil
}
