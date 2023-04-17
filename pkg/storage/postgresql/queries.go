package postgresql

import (
	"database/sql"
	"fmt"
	"math/rand"
	e "telegram-bot/solte.lab/pkg/errhandler"
	"telegram-bot/solte.lab/pkg/storage"
	"time"
)

func (s *Storage) Save(p *storage.Page) (err error) {
	defer func() { err = e.WrapIfErr("can't save page to database", err) }()

	id, err := s.getUserID(p.UserName)
	if err == sql.ErrNoRows {
		id, err = s.InsertUserReturnID(p)
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	_, err = s.db.Exec(`INSERT INTO links (user_id, link) VALUES ($1, $2)`, id, p.URL)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) PickRandom(username string) (page *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("can't pick random page from file", err) }()

	query := `SELECT user_name, links.link, links.id
	FROM users
	JOIN links on links.user_id = users.id
	WHERE user_name = $1;`

	var pages []storage.Page

	rows, err := s.db.Query(query, username)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var p storage.Page
		err = rows.Scan(&p.UserName, &p.URL, &p.URLId)
		if err != nil {
			return nil, err
		}
		pages = append(pages, p)
	}

	if len(pages) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	source := rand.NewSource(time.Now().UnixNano())
	rand.New(source)
	n := rand.Intn(len(pages))

	page = &pages[n]

	return page, nil
}

func (s *Storage) Remove(p *storage.Page) error {
	deleteStmt := `delete from links where id=$1`

	_, err := s.db.Exec(deleteStmt, p.URLId)
	if err != nil {
		return e.Wrap("can't remove file", err)
	}
	return nil
}

func (s *Storage) IsExist(p *storage.Page) (bool, error) {
	id, err := s.getUserID(p.UserName)
	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	p.UserID = id
	return s.checkLink(p)
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

func (s *Storage) InsertUserReturnID(p *storage.Page) (int, error) {
	var userID int
	tx, err := s.db.Begin()

	stmt, err := tx.Prepare(`INSERT INTO "users" (user_name) VALUES ($1) RETURNING id`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(p.UserName).Scan(&userID)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()

	if err != nil {
		return 0, err
	}

	return userID, nil
}
