package query

import (
	"database/sql"

	"github.com/cresendoo/decidash-backend/pkg/errorx"
)

func LockTable(conn Conn, v Model) (sql.Result, error) {
	res, err := conn.Exec("LOCK TABLES " + v.TableName() + " WRITE")
	if err != nil {
		return nil, errorx.Wrap(err)
	}
	return res, nil
}

func UnLockTables(conn Conn) (sql.Result, error) {
	res, err := conn.Exec("UNLOCK TABLES")
	if err != nil {
		return nil, errorx.Wrap(err)
	}
	return res, nil
}
