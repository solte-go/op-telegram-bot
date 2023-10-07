package scheduler

import (
	"context"
	"database/sql"
	"time"

	"go.uber.org/zap"
	"telegram-bot/solte.lab/pkg/config"
	e "telegram-bot/solte.lab/pkg/errhandler"
	"telegram-bot/solte.lab/pkg/events/telegram"
	"telegram-bot/solte.lab/pkg/models"
	"telegram-bot/solte.lab/pkg/queue/kafka"
	"telegram-bot/solte.lab/pkg/storage"
	"telegram-bot/solte.lab/pkg/storage/postgresql"
	"telegram-bot/solte.lab/pkg/storage/repositories/user"
)

type Scheduler struct {
	conf      *config.Scheduler
	store     storage.Repository
	publisher *kafka.Publisher
	logger    *zap.Logger
	repo      *user.Repository
}

func New(db *postgresql.PostgresStorage, publisher *kafka.Publisher, logger *zap.Logger, conf *config.Scheduler) *Scheduler {
	return &Scheduler{
		store:     db,
		publisher: publisher,
		logger:    logger,
		conf:      conf,
		repo:      new(user.Repository),
	}
}

func (s *Scheduler) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.conf.Interval)
	for {
		select {
		case <-ticker.C:
			err := s.interact(ctx)
			if err != nil {
				s.logger.Warn("can't interact with user", zap.Error(err))
			}

		case <-ctx.Done():
			ticker.Stop()
			s.logger.Warn("Scheduler: context is canceled, exiting.")
			return nil
		}
	}
}

func (s *Scheduler) interact(ctx context.Context) error {
	tx, err := s.store.BeginTx(ctx, 0, false)
	if err != nil {
		return err
	}

	users, err := s.repo.GetUsersForInteraction(tx)
	if err != nil {
		return s.store.HandleError(ctx, err, tx)
	}

	for _, u := range users {
		err := s.prepareMessage(ctx, &u, tx)
		if err != nil {
			return err
		}
	}

	err = s.store.CommitTx(ctx, tx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Scheduler) prepareMessage(ctx context.Context, user *models.User, tx *sql.Tx) error {
	if user.NextInteraction.Before(time.Now().UTC()) {
		user.Cmd = telegram.CmdInteract
		user.NextInteraction = time.Now().Add(10 * time.Minute).UTC()

		message, err := s.publisher.PrepareMessage(user)
		if err != nil {
			return e.Wrap("can't prepare message for kafka", err)
		}

		s.publisher.SendMessage(ctx, message)

		err = s.repo.SetUserNextInteractionTime(tx, user)
		if err != nil {
			return s.store.HandleError(ctx, err, tx)
		}
	}

	return nil
}
