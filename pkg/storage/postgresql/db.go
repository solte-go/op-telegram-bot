package postgresql

import (
	"database/sql"
	"fmt"
	"telegram-bot/solte.lab/pkg/storage/postgresql/internal"

	"telegram-bot/solte.lab/pkg/config"
)

//var storagePool map[string]*Storage
//
//func GetStorage(alias string) (*Storage, error) {
//	if storagePool == nil {
//		return nil, fmt.Errorf("storage pool is empty")
//	}
//	storage, ok := storagePool[alias]
//	if ok {
//		return storage, nil
//	}
//	return nil, fmt.Errorf("storage with alias %s not found", alias)
//}

type PostgresStorage struct {
	db *sql.DB
}

func New(conf *config.Postgres) (*PostgresStorage, error) {
	db, err := newDB(conf.OPDB)
	if err != nil {
		return nil, err
	}

	storage := &PostgresStorage{db: db}

	err = storage.init(conf.OPDB.Alias)
	if err != nil {
		return nil, fmt.Errorf("can't initialize storage: %w", err)
	}

	//st := make(map[string]*Storage)
	//storagePool = st
	//storagePool[conf.OPDB.Alias] = storage

	return storage, nil
}

func newDB(conf *config.PostgresSQLConfig) (*sql.DB, error) {
	psqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", conf.Host, conf.Port, conf.Username, conf.Password, conf.DBName)

	db, err := sql.Open("postgres", psqlConn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, err
}

func (s *PostgresStorage) Close() error {
	err := s.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStorage) DropTables() error {
	q := `DROP TABLE IF EXISTS users, words, links CASCADE;`

	_, err := s.db.Exec(q)
	if err != nil {
		return fmt.Errorf("can't drop tables: %w", err)
	}

	return nil
}

func (s *PostgresStorage) init(alias string) error {
	q := internal.CreateTables(alias)
	_, err := s.db.Exec(q)
	if err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}

	fmt.Println("Tables created successfully")

	return nil
}
