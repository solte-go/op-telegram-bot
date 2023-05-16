package postgresql

import (
	"fmt"
	"telegram-bot/solte.lab/pkg/models"
)

func (s *Storage) CreateUser(user *models.Admin) error {
	if err := s.db.QueryRow(
		"INSERT INTO admins (user_name, email, hashed_password, hashed_token) VALUES ($1, $2, $3, $4) RETURNING id",
		user.Name,
		user.Email,
		user.HashedPassword,
		user.HashedToken,
	).Scan(&user.Id); err != nil {
		return err
	}
	return nil
}

func (s *Storage) FindByEmail(email string) (*models.Admin, error) {

	query := `SELECT id, user_name, email, hashed_password, hashed_token FROM admins WHERE email = $1`

	user := &models.Admin{}
	if err := s.db.QueryRow(query, email).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.HashedPassword,
		&user.HashedToken,
	); err != nil {
		return nil, err
	}
	return user, nil
}

//TODO move to server

// SessionSave Perform generation of cookie for user
func (s *Storage) SessionSave(user *models.Admin) error {
	err := s.UpdateUserData(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) UpdateUserData(user *models.Admin) error {
	if err := s.db.QueryRow(
		"UPDATE admins SET (email, hashed_password, hashed_token) = ($1, $2, $3) WHERE email = $4 RETURNING id",
		user.Email,
		user.HashedPassword,
		user.HashedToken,
		user.Email,
	).Scan(&user.Id); err != nil {
		return err
	}
	return nil
}

func (s *Storage) FindBySessionToken(HashedToken string) (*models.Admin, error) {
	user := &models.Admin{}

	query := `SELECT id, user_name, email, hashed_password, hashed_token FROM admins WHERE hashed_token = $1`

	if err := s.db.QueryRow(
		query,
		HashedToken,
	).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.HashedPassword,
		&user.HashedToken,
	); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Storage) AddNewWordsToDataBase(words []models.Words) (err error) {

	query := "INSERT INTO words (suomi, english, russian, letter, topic) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	for _, word := range words {
		_, err = s.db.Exec(
			query,
			word.Suomi,
			word.English,
			word.Russian,
			word.Letter,
			word.Topic,
		)
		if err != nil {
			return fmt.Errorf("error during inserting word %s %w", word.Suomi, err)
		}

	}
	return nil
}
