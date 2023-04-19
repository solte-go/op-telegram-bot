package dialect

import (
	"errors"
	"time"
)

type Dialect struct {
	Alphabet Alphabet
	Topics   Topic
}

type Alphabet struct {
	Letters  []string
	UpdateAt time.Time
}

type Topic struct {
	Titles   []string
	UpdateAt time.Time
}

var ErrUnsupportedTopic = errors.New("unsupported topic argument")

func New() *Dialect {
	return &Dialect{}
}

func (d *Dialect) SyncAlphabet() bool {
	if d.Alphabet.Letters == nil || time.Now().After(d.Alphabet.UpdateAt.Add(+6*time.Hour)) {
		return true
	}
	return false
}

func (d *Dialect) SyncTopics() bool {
	if d.Topics.Titles == nil || time.Now().After(d.Topics.UpdateAt.Add(+6*time.Hour)) {
		return true
	}
	return false
}

func (t *Topic) CheckTopic(topic string) error {
	if topic == "" {
		return ErrUnsupportedTopic
	}

	for _, v := range t.Titles {
		if v == topic {
			return nil
		}
	}
	return ErrUnsupportedTopic
}
