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

type Administrators interface {
	CreateUser(user *models.Admin) error
	FindByEmail(email string) (*models.Admin, error)
	SessionSave(user *models.Admin) error
	UpdateUserData(user *models.Admin) error
	FindBySessionToken(HashedToken string) (*models.Admin, error)
	AddNewWordsToDataBase(words []models.Words) error
}
