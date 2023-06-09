package postgresql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"telegram-bot/solte.lab/pkg/config"
)

type Storage struct {
	db *sql.DB
}

func New(conf *config.PostgreSQL) (*Storage, error) {
	db, err := newDB(conf)
	if err != nil {
		return nil, err
	}

	storage := &Storage{db: db}

	err = storage.init()
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("can't initialize storage: %w", err)
	}

	return storage, nil
}

func newDB(conf *config.PostgreSQL) (*sql.DB, error) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", conf.Host, conf.Port, conf.Username, conf.Password, conf.DBName)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, err
}

func (s *Storage) Close() error {
	err := s.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) DropTables() error {
	q := `DROP TABLE IF EXISTS users, links CASCADE;`

	_, err := s.db.Exec(q)
	if err != nil {
		return fmt.Errorf("can't drop tables: %w", err)
	}

	return nil
}

func (s *Storage) init() error {
	q := `CREATE TABLE IF NOT EXISTS public.links (
    	id SERIAL PRIMARY KEY,
    	user_id integer NOT NULL,
   		link TEXT NOT NULL CONSTRAINT "Links_pk" UNIQUE,
    	create_at timestamp DEFAULT CURRENT_TIMESTAMP
	);

	alter table links
    	owner to postgres;

	CREATE TABLE IF NOT EXISTS public.users (
    	id SERIAL PRIMARY KEY,
    	user_name varchar not null
        constraint "Users_pk"
            unique
	);

	alter table users
    	owner to postgres;
`
	_, err := s.db.Exec(q)
	if err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}

	fmt.Println("Tables created successfully")

	return nil
}
