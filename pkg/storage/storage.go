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

//func (p *models.Page) Hash() (string, error) {
//	h := sha1.New()
//	if _, err := io.WriteString(h, p.URL); err != nil {
//		return "", e.Wrap("can't create hash", err)
//	}
//
//	if _, err := io.WriteString(h, p.UserName); err != nil {
//		return "", e.Wrap("can't create hash", err)
//	}
//	return fmt.Sprintf("%x", h.Sum(nil)), nil
//}
