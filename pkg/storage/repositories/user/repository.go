package user

import (
	"database/sql"
	"time"

	e "telegram-bot/solte.lab/pkg/errhandler"
	"telegram-bot/solte.lab/pkg/models"
)

type ConnectionContract interface {
	Query(query string, args ...any) (*sql.Rows, error)
	Exec(query string, args ...any) (sql.Result, error)
}

type Repository struct {
}

func (s *Repository) GetUsersForInteraction(db ConnectionContract) (users []models.User, err error) {
	defer func() { err = e.WrapIfErr("can't get users from database", err) }()

	query := `SELECT id,user_name, topic, user_language, seq_offset, chat_id, next_interaction_at
			  	FROM users 
			  	WHERE interaction=true 
		      	AND interaction_intensity!=0 
		      	AND chat_id IS NOT NULL;`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var u models.User
		if err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Topic,
			&u.Language,
			&u.Offset,
			&u.ChatID,
			&u.NextInteraction,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (s *Repository) SetUserNextInteractionTime(db ConnectionContract, user *models.User) (err error) {
	defer func() { err = e.WrapIfErr("can't set user parameters", err) }()
	query := `UPDATE users SET next_interaction_at=$2 WHERE id=$1;`

	_, err = db.Exec(query, user.ID, user.NextInteraction.Format(time.RFC3339))
	if err != nil {
		return err
	}

	return err
}
