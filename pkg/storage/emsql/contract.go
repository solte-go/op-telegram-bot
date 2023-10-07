package emsql

import (
	"telegram-bot/solte.lab/pkg/models"
	"telegram-bot/solte.lab/pkg/storage/postgresql"
)

type OpSQL struct {
	st *postgresql.PostgresStorage
	//cache *cache.Container
}

type UserContract interface {
	UserExist(user *models.User) (bool, error)
	GetUser(user *models.User) (err error)
	GetAllUsers() (users []models.User, err error)
	InsertUser(user *models.User) (err error)
	UpdateUserLang(user *models.User) error
	UpdateUserTopic(user *models.User) error
	UpdateUserOffset(user *models.User) error
}

type DialectContract interface {
	GetWords(offset int) (words []*models.Words, newOffset int, err error)
	GetWordsFromTopic(topicTitle string, offset int) (words []*models.Words, newOffset int, err error)
	GetAlphabet() ([]string, error)
	GetTopics() ([]string, error)
}

type CacheContract interface {
	AddUser(user *models.User)
	GetUser(name string) (*models.User, bool)
	UpdateUser(user *models.User) error
	UpdateUserWithUpset(user *models.User) error
}

type AdminsContract interface {
	CreateUser(user *models.Admin) error
	FindByEmail(email string) (*models.Admin, error)
	SessionSave(user *models.Admin) error
	UpdateUserData(user *models.Admin) error
	FindBySessionToken(HashedToken string) (*models.Admin, error)
	AddNewWordsToDataBase(words []models.Words) error
}

type OPContract interface {
	User() UserContract
	Dialect() DialectContract
	Admin() AdminsContract
	//Cache() CacheContract
}

type ServiceContract interface {
	Close() error
	DropTables() error
}

func New(storage *postgresql.PostgresStorage) OPContract {
	return &OpSQL{
		st: storage,
		//cache: cache.New(ctx, storage),
	}
}

func (op *OpSQL) User() UserContract {
	return op.st
}

func (op *OpSQL) Admin() AdminsContract {
	return op.st
}

func (op *OpSQL) Dialect() DialectContract {
	return op.st
}

//func (op *OpSQL) Cache() CacheContract {
//	return op.cache
//}
