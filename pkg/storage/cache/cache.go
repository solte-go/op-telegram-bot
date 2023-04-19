package cache

import (
	"context"
	"errors"
	"math/rand"
	"telegram-bot/solte.lab/pkg/config"
	e "telegram-bot/solte.lab/pkg/errhandler"
	"telegram-bot/solte.lab/pkg/models"
	"telegram-bot/solte.lab/pkg/storage"
	"telegram-bot/solte.lab/pkg/storage/cache/syncContainer"
	"telegram-bot/solte.lab/pkg/storage/dialect"
	"telegram-bot/solte.lab/pkg/storage/postgresql"
	"time"
)

type StorageCache struct {
	dialect *dialect.Dialect
	db      databaseContract
	cache   cacheContainer
}

func New(ctx context.Context, conf *config.PostgreSQL) (*StorageCache, error) {
	d := dialect.New()
	st, err := postgresql.New(conf)
	if err != nil {
		return nil, err
	}

	return &StorageCache{
		dialect: d,
		db:      st,
		cache:   syncContainer.New(ctx, st),
	}, nil
}

type databaseContract interface {
	GetWords(letter string) (page []*storage.Words, err error)
	GetAlphabet() ([]string, error)
	GetUser(user *models.User) (err error)
	InsertUser(user *models.User) (err error)
	UserExist(user *models.User) (bool, error)
	UpdateUserLang(user *models.User) error
}

type cacheContainer interface {
	AddUser(user *models.User)
	GetUser(name string) (models.User, bool)
	UpdateUser(user *models.User) error
}

func (s *StorageCache) SetUserLanguage(user *models.User) (err error) {
	defer func() { err = e.WrapIfErr("can't update user settings", err) }()

	ok, err := s.db.UserExist(user)
	if err != nil {
		return err
	}

	if !ok {
		return s.db.InsertUser(user)
	}

	err = s.cache.UpdateUser(user)
	if err != nil {
		if !errors.Is(err, syncContainer.ErrorNoUserForUpdate) {
			return err
		}
	}

	return s.db.UpdateUserLang(user)
}

func (s *StorageCache) PickRandomWords(user *models.User) (page *storage.Words, err error) {
	if s.dialect.Sync() {
		s.dialect.Alphabet, err = s.db.GetAlphabet()
		if err != nil {
			return nil, err
		}
	}

	cachedUser, ok := s.cache.GetUser(user.Name)
	if ok {
		user.Topic = cachedUser.Topic
		user.Language = cachedUser.Language
	} else {
		err := s.db.GetUser(user)
		if err != nil {
			return nil, err
		}
		s.cache.AddUser(user)
	}

	source := rand.NewSource(time.Now().UnixNano())
	rand.New(source)
	n := rand.Intn(len(s.dialect.Alphabet))

	res, err := s.db.GetWords(s.dialect.Alphabet[n])
	if err != nil {
		return nil, err
	}

	n = rand.Intn(len(res))
	rndWord := res[n]

	return rndWord, nil
}

//func (s *StorageCache) Remove(p *storage.Page) error {
//
//	return nil
//}
//
//func (s *StorageCache) IsExist(p *storage.Page) (bool, error) {
//
//	return s.checkLink(p)
//}

//func (s *StorageCache) Save(p *storage.Page) (err error) {
//	defer func() { err = e.WrapIfErr("can't save page to database", err) }()
//
//	return nil
//}
//
//func (s *StorageCache) PickRandom(username string) (page *storage.Page, err error) {
//	defer func() { err = e.WrapIfErr("can't pick random page from file", err) }()
//
//	return page, nil
//}
