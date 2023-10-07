package storage

import (
	"context"
	"database/sql"

	"telegram-bot/solte.lab/pkg/models"
)

type Storage interface {
	GetUser(user *models.User) (err error)
	InsertUser(user *models.User) (err error)
	UserExist(user *models.User) (bool, error)
	UpdateUserLang(user *models.User) error
	UpdateUserTopic(user *models.User) error
	UpdateUserOffset(user *models.User) error
	GetWords(offset int) (words []*models.Words, newOffset int, err error)
	GetWordsFromTopic(topicTitle string, offset int) (words []*models.Words, newOffset int, err error)
	GetAlphabet() ([]string, error)
	GetTopics() ([]string, error)
	GetAllUsers() (users []models.User, err error)
	GeUsersForInteraction(ctx context.Context) (users []models.User, err error)

	//type AdminsContract interface {
	CreateUser(user *models.Admin) error
	FindByEmail(email string) (*models.Admin, error)
	SessionSave(user *models.Admin) error
	UpdateUserData(user *models.Admin) error
	FindBySessionToken(HashedToken string) (*models.Admin, error)
	AddNewWordsToDataBase(words []models.Words) error
}

type Repository interface {
	BeginTx(ctx context.Context, isolation sql.IsolationLevel, readOnly bool) (*sql.Tx, error)
	HandleError(ctx context.Context, err error, tx *sql.Tx) error
	CommitTx(ctx context.Context, tx *sql.Tx) error
}

//type UserContract interface {
//	GetUser(user *models.User) (err error)
//	InsertUser(user *models.User) (err error)
//	UserExist(user *models.User) (bool, error)
//	UpdateUserLang(user *models.User) error
//	UpdateUserTopic(user *models.User) error
//	UpdateUserOffset(user *models.User) error
//	GetAllUsers() (users []models.User, err error)
//}

//type DialectContract interface {
//	GetWords(offset int) (words []*models.Words, newOffset int, err error)
//	GetWordsFromTopic(topicTitle string, offset int) (words []*models.Words, newOffset int, err error)
//	GetAlphabet() ([]string, error)
//	GetTopics() ([]string, error)
//}
