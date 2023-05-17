package worker

import (
	"context"
	"errors"
	"math/rand"
	e "telegram-bot/solte.lab/pkg/errhandler"
	"telegram-bot/solte.lab/pkg/models"
	"telegram-bot/solte.lab/pkg/storage/cache"
	"telegram-bot/solte.lab/pkg/storage/dialect"
	"telegram-bot/solte.lab/pkg/storage/emsql"
	"time"
)

type Worker struct {
	dialect *dialect.Dialect
	storage emsql.OPContract
	cache   cacheContract
}

type cacheContract interface {
	UpdateUser(user *models.User) error
	UpdateUserWithUpset(user *models.User) error
	GetUser(name string) (*models.User, bool)
	AddUser(user *models.User)
}

func New(ctx context.Context, st emsql.OPContract) *Worker {
	d := dialect.New()
	return &Worker{
		dialect: d,
		storage: st,
		cache:   cache.New(ctx, st),
	}
}

func (w *Worker) SetUserLanguage(user *models.User) (err error) {
	defer func() { err = e.WrapIfErr("can't update user settings", err) }()

	ok, err := w.storage.User().UserExist(user)
	if err != nil {
		return err
	}

	if !ok {
		return w.storage.User().InsertUser(user)
	}

	err = w.cache.UpdateUserWithUpset(user)
	if err != nil {
		return err
	}

	return w.storage.User().UpdateUserLang(user)
}

func (w *Worker) GetTopics() (topics []string, err error) {
	if w.dialect.SyncTopics() {
		w.dialect.Topics.Titles, err = w.storage.Dialect().GetTopics()
		if err != nil {
			return nil, err
		}
	}

	topics = w.dialect.Topics.Titles

	return topics, nil
}

func (w *Worker) SetUserTopic(user *models.User, topic string) (err error) {
	defer func() { err = e.WrapIfErr("can't update user settings", err) }()

	if err = w.dialect.Topics.CheckTopic(topic); err != nil {
		return err
	}

	user.Topic = topic
	user.Offset = 0

	ok, err := w.storage.User().UserExist(user)
	if err != nil {
		return err
	}

	if !ok {
		return w.storage.User().InsertUser(user)
	}

	err = w.cache.UpdateUserWithUpset(user)
	if err != nil {
		return err
	}

	return w.storage.User().UpdateUserTopic(user)
}

func (w *Worker) PickRandomWord(user *models.User) (word *models.Words, err error) {
	defer func() { err = e.WrapIfErr("wrapper can't process request", err) }()

	var words []*models.Words
	var newOffset int

	if w.dialect.SyncTopics() {
		w.dialect.Topics.Titles, err = w.storage.Dialect().GetTopics()
		if err != nil {
			return nil, err
		}
	}

	cachedUser, ok := w.cache.GetUser(user.Name)
	if ok {
		user.Topic = cachedUser.Topic
		user.Language = cachedUser.Language
		user.Offset = cachedUser.Offset
		user.Sequence = cachedUser.Sequence
	} else {
		err = w.storage.User().GetUser(user)
		if err != nil {
			return nil, err
		}
		w.cache.AddUser(user)
	}

	if user.Sequence.Words == nil || len(user.Sequence.Words) == 0 || user.Sequence.NeedToUpdate() {
		if user.Topic != "" && !user.IsTopicDefault() {
			words, newOffset, err = w.storage.Dialect().GetWordsFromTopic(user.Topic, user.Offset)
			if err != nil {
				return nil, err
			}
		} else {
			words, newOffset, err = w.storage.Dialect().GetWords(user.Offset)
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

		err = w.storage.User().UpdateUserOffset(user)
		if err != nil {
			return nil, err
		}
	}

	word, update := user.Sequence.GetNextWord()
	if update {
		return nil, errors.New("can't get next word")
	}

	err = w.cache.UpdateUser(user)
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
