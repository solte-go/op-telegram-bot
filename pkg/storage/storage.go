package storage

import (
	"telegram-bot/solte.lab/pkg/models"
)

type Storage interface {
	GetTopics() (topics []string, err error)
	PickRandomWord(user *models.User) (word *models.Words, err error)
	SetUserLanguage(user *models.User) (err error)
	SetUserTopic(user *models.User, topic string) (err error)
}
