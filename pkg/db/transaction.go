package db

import (
	"database/sql"

	"github.com/cresendoo/decidash-backend/pkg/errorx"
)

type Beginner interface {
	Begin() (*sql.Tx, error)
}

type TxFunc func(tx *sql.Tx) error

func Transaction(conn Beginner, txFunc TxFunc) error {
	tx, err := conn.Begin()
	if err != nil {
		return errorx.Wrap(err)
	}
	if err := txFunc(tx); err != nil {
		_ = tx.Rollback()
		return errorx.Wrap(err)
	}

	return errorx.Wrap(tx.Commit())
}
