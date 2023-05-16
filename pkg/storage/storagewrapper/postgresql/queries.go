package postgresql

import (
	"fmt"
	e "telegram-bot/solte.lab/pkg/errhandler"
	"telegram-bot/solte.lab/pkg/models"
	"telegram-bot/solte.lab/pkg/storage"
)

const offsetInc = 20

func (s *Storage) GetWords(offset int) (words []*models.Words, newOffset int, err error) {
	defer func() { err = e.WrapIfErr("can't retrieve words from database", err) }()

	wordIndex, err := s.getOffset()
	if err != nil {
		return words, 0, err
	}

	if offset >= wordIndex {
		offset = 0
	}

	query := `SELECT topic, suomi, russian, english FROM words OFFSET $1 LIMIT $2`

	rows, err := s.db.Query(query, offset, offsetInc)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var word models.Words
		err = rows.Scan(&word.Topic, &word.Suomi, &word.Russian, &word.English)
		if err != nil {
			return nil, 0, err
		}
		words = append(words, &word)
	}

	if len(words) == 0 {
		return nil, 0, storage.ErrNoSavedPages
	}

	newOffset = offset + len(words)

	return words, newOffset, err
}

func (s *Storage) GetWordsFromTopic(topicTitle string, offset int) (words []*models.Words, newOffset int, err error) {
	defer func() { err = e.WrapIfErr("can't retrieve words with specific topic from database", err) }()

	wordIndex, err := s.getOffsetWithTopic(topicTitle)
	if err != nil {
		return words, 0, err
	}

	if offset >= wordIndex {
		offset = 0
	}

	query := `SELECT topic, suomi, russian, english FROM words WHERE topic=$1 OFFSET $2 LIMIT $3`

	rows, err := s.db.Query(query, topicTitle, offset, offsetInc)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var word models.Words
		err = rows.Scan(&word.Topic, &word.Suomi, &word.Russian, &word.English)
		if err != nil {
			return nil, 0, err
		}
		words = append(words, &word)
	}

	if len(words) == 0 {
		return nil, 0, storage.ErrNoSavedPages
	}

	newOffset = offset + len(words)

	return words, newOffset, err
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

func (s *Storage) Remove(p *models.Page) error {
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

func (s *Storage) getOffset() (int, error) {
	var count int

	q := `SELECT COUNT(*) FROM words;`

	if err := s.db.QueryRow(q).Scan(&count); err != nil {
		return 0, fmt.Errorf("can't check if page exists: %w", err)
	}

	if count == 0 {
		return 0, fmt.Errorf("can't get offset")
	}

	return count, nil
}

func (s *Storage) getOffsetWithTopic(topic string) (int, error) {
	var count int

	q := `SELECT COUNT(*) FROM words WHERE topic=$1;`

	if err := s.db.QueryRow(q, topic).Scan(&count); err != nil {
		return 0, fmt.Errorf("can't check if page exists: %w", err)
	}

	if count == 0 {
		return 0, fmt.Errorf("can't get offset")
	}

	return count, nil
}

func (s *Storage) checkLink(p *models.Page) (bool, error) {
	var count int

	q := `SELECT COUNT(*) FROM links WHERE user_id = $1 AND link = $2`

	if err := s.db.QueryRow(q, p.UserID, p.URL).Scan(&count); err != nil {
		return false, fmt.Errorf("can't check if page exists: %w", err)
	}

	return count > 0, nil
}
