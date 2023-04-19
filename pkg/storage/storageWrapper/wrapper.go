package storageWrapper

import (
	"context"
	"errors"
	"math/rand"
	"telegram-bot/solte.lab/pkg/config"
	e "telegram-bot/solte.lab/pkg/errhandler"
	"telegram-bot/solte.lab/pkg/models"
	"telegram-bot/solte.lab/pkg/storage"
	"telegram-bot/solte.lab/pkg/storage/dialect"
	"telegram-bot/solte.lab/pkg/storage/storageWrapper/cache"
	"telegram-bot/solte.lab/pkg/storage/storageWrapper/postgresql"
	"time"
)

type StorageCache struct {
	dialect *dialect.Dialect
	storage *contract
}

func New(ctx context.Context, conf *config.PostgreSQL) (*StorageCache, error) {
	d := dialect.New()
	st, err := postgresql.New(conf)
	if err != nil {
		return nil, err
	}

	return &StorageCache{
		dialect: d,
		storage: newContract(ctx, st),
	}, nil
}

type contract struct {
	user    userContract
	dialect dialectContract
	cache   cacheContract
}

func newContract(ctx context.Context, st *postgresql.Storage) *contract {
	return &contract{
		user:    st,
		dialect: st,
		cache:   cache.New(ctx, st),
	}
}

type userContract interface {
	GetUser(user *models.User) (err error)
	InsertUser(user *models.User) (err error)
	UserExist(user *models.User) (bool, error)
	UpdateUserLang(user *models.User) error
	UpdateUserTopic(user *models.User) error
}

type dialectContract interface {
	GetWords(letter string) (page []*storage.Words, err error)
	GetWordsFromTopic(topicTitle string) (words []*storage.Words, err error)
	GetAlphabet() ([]string, error)
	GetTopics() ([]string, error)
}

type cacheContract interface {
	AddUser(user *models.User)
	GetUser(name string) (models.User, bool)
	UpdateUser(user *models.User) error
}

func (s *StorageCache) SetUserLanguage(user *models.User) (err error) {
	defer func() { err = e.WrapIfErr("can't update user settings", err) }()

	ok, err := s.storage.user.UserExist(user)
	if err != nil {
		return err
	}

	if !ok {
		return s.storage.user.InsertUser(user)
	}

	err = s.storage.cache.UpdateUser(user)
	if err != nil {
		if !errors.Is(err, cache.ErrorNoUserForUpdate) {
			return err
		}
	}

	return s.storage.user.UpdateUserLang(user)
}

func (s *StorageCache) GetTopics() (topics []string, err error) {
	if s.dialect.SyncTopics() {
		s.dialect.Topics.Titles, err = s.storage.dialect.GetTopics()
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

	ok, err := s.storage.user.UserExist(user)
	if err != nil {
		return err
	}

	if !ok {
		return s.storage.user.InsertUser(user)
	}

	err = s.storage.cache.UpdateUser(user)
	if err != nil {
		if !errors.Is(err, cache.ErrorNoUserForUpdate) {
			return err
		}
	}

	return s.storage.user.UpdateUserTopic(user)
}

func (s *StorageCache) PickRandomWord(user *models.User) (page *storage.Words, err error) {
	defer func() { err = e.WrapIfErr("can't get topics form db", err) }()

	if s.dialect.SyncTopics() {
		s.dialect.Topics.Titles, err = s.storage.dialect.GetTopics()
		if err != nil {
			return nil, err
		}
	}

	cachedUser, ok := s.storage.cache.GetUser(user.Name)
	if ok {
		user.Topic = cachedUser.Topic
		user.Language = cachedUser.Language
	} else {
		err = s.storage.user.GetUser(user)
		if err != nil {
			return nil, err
		}
		s.storage.cache.AddUser(user)
	}

	if user.Topic != "" && !user.IsTopicDefault() {
		var res []*storage.Words
		res, err = s.storage.dialect.GetWordsFromTopic(user.Topic)
		if err != nil {
			return nil, err
		}

		source := rand.NewSource(time.Now().UnixNano())
		rand.New(source)

		n := rand.Intn(len(res))
		rndWord := res[n]
		return rndWord, nil
	}

	return s.randomWordsWithoutTopic()
}

func (s *StorageCache) randomWordsWithoutTopic() (page *storage.Words, err error) {
	if s.dialect.SyncAlphabet() {
		s.dialect.Alphabet.Letters, err = s.storage.dialect.GetAlphabet()
		if err != nil {
			return nil, err
		}
	}

	source := rand.NewSource(time.Now().UnixNano())
	rand.New(source)
	n := rand.Intn(len(s.dialect.Alphabet.Letters))

	res, err := s.storage.dialect.GetWords(s.dialect.Alphabet.Letters[n])
	if err != nil {
		return nil, err
	}

	n = rand.Intn(len(res))
	rndWord := res[n]

	return rndWord, nil
}
