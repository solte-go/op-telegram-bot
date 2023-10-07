package models

import (
	"errors"
	"time"
)

const (
	defaultLanguage = "ru"
	defaultTopic    = "all"
	ru              = "ru"
	en              = "en"
)

var ErrUnsupportedLanguage = errors.New("unsupported language argument")

type User struct {
	ID              int       `json:"-"`
	Cmd             string    `json:"cmd"`
	Name            string    `json:"name"`
	Language        string    `json:"language"`
	Topic           string    `json:"topic"`
	ChatID          int       `json:"chat_id"`
	Offset          int       `json:"offset"`
	Sequence        *Sequence `json:"sequence"`
	NextInteraction time.Time
}

func (u *User) SetDefaults() {
	if u.Language == "" {
		u.Language = defaultLanguage
	}
	if u.Topic == "" {
		u.Topic = defaultTopic
	}
}

func (u *User) IsLanguageEnglish() bool {
	if u.Language == en {
		return true
	}
	return false
}

func (u *User) IsTopicDefault() bool {
	if u.Topic == defaultTopic {
		return true
	}
	return false
}

func (u *User) CheckLanguage(lang string) error {
	switch lang {
	case en:
		u.Language = en
	case ru:
		u.Language = ru
	default:
		return ErrUnsupportedLanguage
	}

	return nil
}
