package postgresql

import (
	"database/sql"
	"fmt"

	"golang.org/x/net/context"
	"telegram-bot/solte.lab/pkg/storage/postgresql/internal"

	"telegram-bot/solte.lab/pkg/config"
)

type TransactionHandler interface {
	TransactionError(ctx context.Context, err error, tx *sql.Tx) error
}

type PostgresStorage struct {
	db      *sql.DB
	Handler TransactionHandler
}

func New(conf *config.Postgres) (*PostgresStorage, error) {
	db, err := newDB(conf.OPDB)
	if err != nil {
		return nil, err
	}

	storage := &PostgresStorage{db: db,
		Handler: new(TxErrorHandler),
	}

	//err = storage.init(conf.OPDB.Alias)
	//if err != nil {
	//	return nil, fmt.Errorf("can't initialize storage: %w", err)
	//}

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

func (s *PostgresStorage) BeginTx(ctx context.Context, isolation sql.IsolationLevel, readOnly bool) (*sql.Tx, error) {
	return s.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: isolation,
		ReadOnly:  readOnly,
	})
}

func (s *PostgresStorage) HandleError(ctx context.Context, err error, tx *sql.Tx) error {
	return s.Handler.TransactionError(ctx, err, tx)
}

func (s *PostgresStorage) CommitTx(ctx context.Context, tx *sql.Tx) error {
	err := tx.Commit()
	if err != nil {
		return s.Handler.TransactionError(ctx, fmt.Errorf("postgres: unable to commit transaction: %v", err), tx)
	}
	return nil
}

// TODO move to migration
func (s *PostgresStorage) init(alias string) error {
	q := internal.CreateTables(alias)
	_, err := s.db.Exec(q)
	if err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}

	fmt.Println("Tables created successfully")

	return nil
}
