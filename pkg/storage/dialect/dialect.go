package dialect

import "time"

type Dialect struct {
	Alphabet []string
	UpdateAt time.Time
}

func New() *Dialect {
	return &Dialect{}
}

func (d *Dialect) Sync() bool {
	if d.Alphabet == nil || time.Now().After(d.UpdateAt.Add(+6*time.Hour)) {
		return true
	}
	return false
}
