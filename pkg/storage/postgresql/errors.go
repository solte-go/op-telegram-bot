package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"

	e "telegram-bot/solte.lab/pkg/errorshandler"
)

var reDuplicateKey = regexp.MustCompile(`duplicate key value violates unique constraint`)

func duplicateKeyError(err error) bool {
	return reDuplicateKey.MatchString(err.Error())
}

var reg23505 = regexp.MustCompile(`SQLSTATE 23505`)
var reg23503 = regexp.MustCompile(`SQLSTATE 23503`)

type TxErrorHandler struct {
}

func (t *TxErrorHandler) TransactionError(ctx context.Context, err error, tx *sql.Tx) error {
	//tracing.LogError(ctx, err)

	if tx != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return e.WrapError(fmt.Errorf("postgres: unable to rollback: %v", rollbackErr), e.ErrBadRequest)
		}
	}

	if errors.Is(sql.ErrNoRows, err) {
		return e.WrapError(fmt.Errorf("postgres: %v", err), e.ErrNotFound)
	}

	switch {
	case reg23505.MatchString(err.Error()):
		return e.WrapError(fmt.Errorf("postgres: dublicate key %v", err), e.ErrDuplicateData)

	case reg23503.MatchString(err.Error()):
		return e.WrapError(fmt.Errorf("postgres: foreign key violation %v", err), e.ErrDuplicateData)
	}

	return e.WrapError(fmt.Errorf("postgres: unable to execute query: %v", err), e.ErrBadRequest)
}
