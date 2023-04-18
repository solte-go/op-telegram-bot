package models

const (
	defaultLanguage = "ru"
	defaultTopic    = "all"
	ru              = "ru"
	en              = "en"
)

type User struct {
	ID       int
	Name     string
	Language string
	Topic    string
	ChatID   int
}

func (u *User) SetDefaults() {
	u.Language = defaultLanguage
	u.Topic = defaultTopic
}

func (u *User) IsLanguageEnglish() bool {
	if u.Language == en {
		return true
	}
	return false
}
