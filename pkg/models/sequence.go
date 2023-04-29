package models

type Sequence struct {
	Words    []*Words
	NextWord int
}

func (u *Sequence) ResetSequence() {
	u.NextWord = 0
}

func (u *Sequence) NeedToUpdate() bool {
	if u.NextWord == len(u.Words) {
		return true
	}
	return false
}

func (u *Sequence) GetNextWord() (*Words, bool) {
	if u.NextWord == len(u.Words) {
		u.NextWord = 0
		return nil, true
	}

	word := u.Words[u.NextWord]
	u.NextWord++
	return word, false
}
