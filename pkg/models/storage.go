package models

import (
	"crypto/sha1"
	"fmt"
	"io"
	e "telegram-bot/solte.lab/pkg/errhandler"
)

type Page struct {
	UserID   int
	UserName string
	URLId    int
	URL      string
}

type Words struct {
	Topic   string
	Letter  string
	Suomi   string
	Russian string
	English string
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
