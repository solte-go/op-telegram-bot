package postgresql

import (
	"database/sql"
	"fmt"
	"telegram-bot/solte.lab/pkg/models"

	e "telegram-bot/solte.lab/pkg/errhandler"
)

func (s *Storage) GetAllUsers() (users []models.User, err error) {
	defer func() { err = e.WrapIfErr("can't get users from database", err) }()

	query := `SELECT user_name, topic, user_language, seq_offset FROM users`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.Name, &u.Topic, &u.Language, &u.Offset); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (s *Storage) GetUser(user *models.User) (err error) {
	defer func() { err = e.WrapIfErr("can't get user from database", err) }()

	query := `SELECT topic, user_language
	FROM users
	WHERE user_name = $1;`

	err = s.db.QueryRow(query, user.Name).Scan(&user.Topic, &user.Language)
	if err == sql.ErrNoRows {
		err = s.InsertUser(user)
		if err != nil {
			return err
		}
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) InsertUser(user *models.User) (err error) {
	defer func() { err = e.WrapIfErr("can't insert user to database", err) }()

	query := `INSERT INTO users (user_name, topic, user_language) values ($1, $2, $3) RETURNING id;`

	user.SetDefaults()

	_, err = s.db.Exec(query, user.Name, user.Topic, user.Language)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdateUserLang(user *models.User) error {
	query := `UPDATE users SET user_language = $1 WHERE user_name = $2;`

	_, err := s.db.Exec(query, user.Language, user.Name)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdateUserOffset(user *models.User) error {
	query := `UPDATE users SET seq_offset = $1 WHERE user_name = $2;`

	_, err := s.db.Exec(query, user.Offset, user.Name)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdateUserTopic(user *models.User) error {
	query := `UPDATE users SET topic = $1 WHERE user_name = $2;`

	_, err := s.db.Exec(query, user.Topic, user.Name)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) UserExist(user *models.User) (bool, error) {
	var count int

	query := `SELECT COUNT(*) 
	FROM users WHERE user_name = $1`

	if err := s.db.QueryRow(query, user.Name).Scan(&count); err != nil {
		return false, fmt.Errorf("can't check if user exists: %w", err)
	}

	return count > 0, nil
}

func (s *Storage) insertNewUserReturnID(username string) (int, error) {
	var userID int
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}

	stmt, err := tx.Prepare(`INSERT INTO "users" (user_name) VALUES ($1) RETURNING id`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(username).Scan(&userID)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()

	if err != nil {
		return 0, err
	}

	return userID, nil
}
