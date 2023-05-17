package telegram

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"telegram-bot/solte.lab/pkg/clients/telegram"
	e "telegram-bot/solte.lab/pkg/errhandler"
	"telegram-bot/solte.lab/pkg/events"
	"telegram-bot/solte.lab/pkg/models"
	"telegram-bot/solte.lab/pkg/queue"
	"telegram-bot/solte.lab/pkg/queue/kafka"
)

type Processor struct {
	tg       *telegram.Client
	producer queue.Producer
	logger   *zap.Logger
	offset   int
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *telegram.Client, kafka *kafka.Publisher, logger *zap.Logger) *Processor {
	return &Processor{
		tg:       client,
		producer: kafka,
		logger:   logger,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}
	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}
	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	user := &models.User{
		Name:   meta.Username,
		ChatID: meta.ChatID,
		Cmd:    event.Text,
		Sequence: &models.Sequence{
			Words:    nil,
			NextWord: 0,
		},
	}

	message, err := p.producer.PrepareMessage(user)
	if err != nil {
		return e.Wrap("can't prepare message for kafka", err)
	}

	p.producer.SendMessage(context.TODO(), message)

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}
	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}
