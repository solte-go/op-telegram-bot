package postgresql

import (
	"fmt"
	e "telegram-bot/solte.lab/pkg/errhandler"
	"telegram-bot/solte.lab/pkg/storage"
)

func (s *Storage) GetWords(letter string) (words []*storage.Words, err error) {

	query := `SELECT topic, suomi, russian, english FROM words WHERE letter=$1`

	rows, err := s.db.Query(query, letter)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var word storage.Words
		err = rows.Scan(&word.Topic, &word.Suomi, &word.Russian, &word.English)
		if err != nil {
			return nil, err
		}
		words = append(words, &word)
	}

	if len(words) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	return words, err
}

func (s *Storage) GetWordsFromTopic(topicTitle string) (words []*storage.Words, err error) {

	query := `SELECT topic, suomi, russian, english FROM words WHERE topic=$1`

	rows, err := s.db.Query(query, topicTitle)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var word storage.Words
		err = rows.Scan(&word.Topic, &word.Suomi, &word.Russian, &word.English)
		if err != nil {
			return nil, err
		}
		words = append(words, &word)
	}

	if len(words) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	return words, err
}

func (s *Storage) GetTopics() ([]string, error) {

	var topics []string
	query := `SELECT DISTINCT topic FROM words`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var title string
		err = rows.Scan(&title)
		if err != nil {
			return nil, err
		}

		topics = append(topics, title)
	}

	topics = append(topics, "all")

	return topics, nil
}

func (s *Storage) Remove(p *storage.Page) error {
	deleteStmt := `delete from links where id=$1`

	_, err := s.db.Exec(deleteStmt, p.URLId)
	if err != nil {
		return e.Wrap("can't remove file", err)
	}
	return nil
}

func (s *Storage) GetAlphabet() ([]string, error) {

	var alphabet []string
	query := `SELECT DISTINCT letter FROM words`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var letter string
		err = rows.Scan(&letter)
		if err != nil {
			return nil, err
		}
		alphabet = append(alphabet, letter)
	}
	return alphabet, nil
}

func (s *Storage) getUserID(userName string) (id int, err error) {
	query := `SELECT id FROM users WHERE user_name = $1`

	err = s.db.QueryRow(query, userName).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Storage) checkLink(p *storage.Page) (bool, error) {
	var count int

	q := `SELECT COUNT(*) FROM links WHERE user_id = $1 AND link = $2`

	if err := s.db.QueryRow(q, p.UserID, p.URL).Scan(&count); err != nil {
		return false, fmt.Errorf("can't check if page exists: %w", err)
	}

	return count > 0, nil
}
