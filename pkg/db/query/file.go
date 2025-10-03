package query

import (
	"bytes"
	"database/sql"
	"os"

	"github.com/cresendoo/decidash-backend/pkg/errorx"
)

func ExecFile(conn Conn, path string) (sql.Result, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errorx.Wrap(err)
	}
	result, err := conn.Exec(bytes.NewBuffer(data).String())
	if err != nil {
		return nil, errorx.Wrap(err)
	}
	return result, nil
}
