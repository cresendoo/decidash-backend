package query

import (
	"database/sql"

	"github.com/cresendoo/decidash-backend/pkg/errorx"
	"github.com/go-sql-driver/mysql"
)

func IsNoRows(err error) bool {
	return errorx.Is(err, sql.ErrNoRows)
}

func IsMySqlError(err error) bool {
	var sqlErr *mysql.MySQLError
	return errorx.As(err, &sqlErr)
}

func IsDuplicateEntry(err error) bool {
	var sqlErr *mysql.MySQLError
	return errorx.As(err, &sqlErr) && sqlErr.Number == 1062
}
