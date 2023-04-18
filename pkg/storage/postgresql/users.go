package postgresql

import (
	"database/sql"
	e "telegram-bot/solte.lab/pkg/errhandler"
	"telegram-bot/solte.lab/pkg/models"
)

func (s *Storage) GetUser(user *models.User) (err error) {
	defer func() { err = e.WrapIfErr("can't get user from database", err) }()

	query := `SELECT topic, user_language
	FROM users
	WHERE user_name = $1;`

	err = s.db.QueryRow(query, user.Name).Scan(&user.Topic, &user.Language)
	if err == sql.ErrNoRows {
		_, err = s.insertNewUserReturnID(user.Name)
		if err != nil {
			return err
		}
		user.SetDefaults()
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) insertNewUserReturnID(username string) (int, error) {
	var userID int
	tx, err := s.db.Begin()

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
