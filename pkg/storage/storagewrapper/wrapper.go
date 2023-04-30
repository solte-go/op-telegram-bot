package storagewrapper

import (
	"context"
	"errors"
	"math/rand"
	"telegram-bot/solte.lab/pkg/config"
	e "telegram-bot/solte.lab/pkg/errhandler"
	"telegram-bot/solte.lab/pkg/models"
	"telegram-bot/solte.lab/pkg/storage/dialect"
	"telegram-bot/solte.lab/pkg/storage/storagewrapper/cache"
	"telegram-bot/solte.lab/pkg/storage/storagewrapper/postgresql"
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
	UpdateUserOffset(user *models.User) error
}

type dialectContract interface {
	GetWords(offset int) (words []*models.Words, newOffset int, err error)
	GetWordsFromTopic(topicTitle string, offset int) (words []*models.Words, newOffset int, err error)
	GetAlphabet() ([]string, error)
	GetTopics() ([]string, error)
}

type cacheContract interface {
	AddUser(user *models.User)
	GetUser(name string) (*models.User, bool)
	UpdateUser(user *models.User) error
	UpdateUserWithUpset(user *models.User) error
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

	err = s.storage.cache.UpdateUserWithUpset(user)
	if err != nil {
		return err
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
	user.Offset = 0

	ok, err := s.storage.user.UserExist(user)
	if err != nil {
		return err
	}

	if !ok {
		return s.storage.user.InsertUser(user)
	}

	err = s.storage.cache.UpdateUserWithUpset(user)
	if err != nil {
		return err
	}

	return s.storage.user.UpdateUserTopic(user)
}

func (s *StorageCache) PickRandomWord(user *models.User) (word *models.Words, err error) {
	defer func() { err = e.WrapIfErr("wrapper can't process request", err) }()

	var words []*models.Words
	var newOffset int

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
		user.Offset = cachedUser.Offset
		user.Sequence = cachedUser.Sequence
	} else {
		err = s.storage.user.GetUser(user)
		if err != nil {
			return nil, err
		}
		s.storage.cache.AddUser(user)
	}

	if user.Sequence.Words == nil || len(user.Sequence.Words) == 0 || user.Sequence.NeedToUpdate() {
		if user.Topic != "" && !user.IsTopicDefault() {
			words, newOffset, err = s.storage.dialect.GetWordsFromTopic(user.Topic, user.Offset)
			if err != nil {
				return nil, err
			}
		} else {
			words, newOffset, err = s.storage.dialect.GetWords(user.Offset)
			if err != nil {
				return nil, err
			}
		}

		user.Offset = newOffset

		// Shuffle words
		source := rand.NewSource(time.Now().UnixNano())
		rand.New(source).Shuffle(len(words), func(i, j int) {
			words[i], words[j] = words[j], words[i]
		})

		user.Sequence.Words = words
		user.Sequence.ResetSequence()

		err = s.storage.user.UpdateUserOffset(user)
		if err != nil {
			return nil, err
		}
	}

	word, update := user.Sequence.GetNextWord()
	if update {
		return nil, errors.New("can't get next word")
	}

	err = s.storage.cache.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	return word, nil
}

//func (s *StorageCache) randomWordsWithoutTopic() (page *models.Words, err error) {
//	if s.dialect.SyncAlphabet() {
//		s.dialect.Alphabet.Letters, err = s.storage.dialect.GetAlphabet()
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	if len(s.dialect.Alphabet.Letters) == 0 {
//		return nil, errors.New("no words in db")
//	}
//
//	source := rand.NewSource(time.Now().UnixNano())
//	rand.New(source)
//	n := rand.Intn(len(s.dialect.Alphabet.Letters))
//
//	res, err := s.storage.dialect.GetWords(s.dialect.Alphabet.Letters[n])
//	if err != nil {
//		return nil, err
//	}
//
//	n = rand.Intn(len(res))
//	rndWord := res[n]
//
//	return page, nil
//}
