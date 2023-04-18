package storage

import (
	"crypto/sha1"
	"fmt"
	"io"
	e "telegram-bot/solte.lab/pkg/errhandler"
	"telegram-bot/solte.lab/pkg/models"
)

//Save(p *Page) error
//PickRandom(username string) (page *Page, err error)
//Remove(p *Page) error
//IsExist(p *Page) (bool, error)

type Storage interface {
	PickRandomWords(user *models.User) (page *Words, err error)
}

func (p *Page) Hash() (string, error) {
	h := sha1.New()
	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("can't create hash", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("can't create hash", err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
