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
	GetWordsFromTopic(topicTitle string) (words []*storage.Words, err error)
	GetAlphabet() ([]string, error)
	GetTopics() ([]string, error)
	GetUser(user *models.User) (err error)
	InsertUser(user *models.User) (err error)
	UserExist(user *models.User) (bool, error)
	UpdateUserLang(user *models.User) error
	UpdateUserTopic(user *models.User) error
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

func (s *StorageCache) GetTopics() (topics []string, err error) {
	if s.dialect.SyncTopics() {
		s.dialect.Topics.Titles, err = s.db.GetTopics()
		if err != nil {
			return nil, err
		}
	}

	topics = s.dialect.Topics.Titles

	return topics, nil
}

func (s *StorageCache) SetUserTopic(user *models.User, topic string) (err error) {
	defer func() { err = e.WrapIfErr("can't update user settings", err) }()

	if err = s.dialect.Topics.CheckTopic(topic); err != nil {
		return err
	}

	user.Topic = topic

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

	return s.db.UpdateUserTopic(user)
}

func (s *StorageCache) PickRandomFromTopic(user *models.User) (page *storage.Words, err error) {
	defer func() { err = e.WrapIfErr("can't get topics form db", err) }()

	if s.dialect.SyncTopics() {
		s.dialect.Topics.Titles, err = s.db.GetTopics()
		if err != nil {
			return nil, err
		}
	}

	cachedUser, ok := s.cache.GetUser(user.Name)
	if ok {
		user.Topic = cachedUser.Topic
		user.Language = cachedUser.Language
	} else {
		err = s.db.GetUser(user)
		if err != nil {
			return nil, err
		}
		s.cache.AddUser(user)
	}

	if user.Topic != "" && !user.IsTopicDefault() {
		var res []*storage.Words
		res, err = s.db.GetWordsFromTopic(user.Topic)
		if err != nil {
			return nil, err
		}

		source := rand.NewSource(time.Now().UnixNano())
		rand.New(source)

		n := rand.Intn(len(res))
		rndWord := res[n]
		return rndWord, nil
	}

	return s.pickRandomWords(user)
}

func (s *StorageCache) pickRandomWords(user *models.User) (page *storage.Words, err error) {
	if s.dialect.SyncAlphabet() {
		s.dialect.Alphabet.Letters, err = s.db.GetAlphabet()
		if err != nil {
			return nil, err
		}
	}

	cachedUser, ok := s.cache.GetUser(user.Name)
	if ok {
		user.Topic = cachedUser.Topic
		user.Language = cachedUser.Language
	} else {
		err = s.db.GetUser(user)
		if err != nil {
			return nil, err
		}
		s.cache.AddUser(user)
	}

	source := rand.NewSource(time.Now().UnixNano())
	rand.New(source)
	n := rand.Intn(len(s.dialect.Alphabet.Letters))

	res, err := s.db.GetWords(s.dialect.Alphabet.Letters[n])
	if err != nil {
		return nil, err
	}

	n = rand.Intn(len(res))
	rndWord := res[n]

	return rndWord, nil
}
