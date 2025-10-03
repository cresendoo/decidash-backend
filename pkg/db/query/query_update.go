package query

import (
	"fmt"
	"log/slog"

	"github.com/cresendoo/decidash-backend/pkg/errorx"
)

func SetStruct(db Conn, v Model) error {
	cols := getStructColumns(v)
	baseQuery, cols := makeQueryString(v.TableName(), cols, stmtUpdate)

	stmt, err := db.Prepare(baseQuery)
	if err != nil {
		return errorx.Wrap(err)
	}
	defer stmt.Close()

	slog.Default().With("query", "UPDATE").Debug(baseQuery)

	if _, err = stmt.Exec(getStructValues(cols, v)...); err != nil {
		return errorx.Wrap(err)
	}
	return nil
}

func SetStructWithID(db Conn, v Model) error {
	cols := getStructColumns(v)
	baseQuery, cols := makeQueryString(v.TableName(), cols, stmtUpdateByID)

	stmt, err := db.Prepare(baseQuery)
	if err != nil {
		return errorx.Wrap(err)
	}
	defer stmt.Close()

	slog.Default().With("query", "UPDATE").Debug(baseQuery)

	if _, err = stmt.Exec(getStructValues(cols, v)...); err != nil {
		return errorx.Wrap(err)
	}
	return nil
}

func SetStructWithCols(db Conn, v Model, updateCols ...string) error {
	allCols := getStructColumns(v)
	var cols []string
	// add primary key
	cols = append(cols, allCols[0])
	cols = append(cols, updateCols...)
	baseQuery, cols := makeQueryString(v.TableName(), cols, stmtUpdateByID)

	stmt, err := db.Prepare(baseQuery)
	if err != nil {
		return errorx.Wrap(err)
	}
	defer stmt.Close()

	slog.Default().With("query", "UPDATE").Debug(baseQuery)

	if _, err = stmt.Exec(getStructValues(cols, v)...); err != nil {
		return errorx.Wrap(err)
	}
	return nil
}

func SetStructsWithCols(db Conn, where *QueryString, v Model, updateCols ...string) error {
	baseQuery, _ := makeQueryString(v.TableName(), updateCols, stmtUpdate)

	var query string
	vals := getStructValues(updateCols, v)
	if where == nil {
		query = baseQuery
	} else {
		query = fmt.Sprintf("%s %s %s", where.PreStmt, baseQuery, where.Stmt)
		vals = append(vals, where.Vals...)
	}

	stmt, err := db.Prepare(query)
	if err != nil {
		return errorx.Wrap(err)
	}
	defer stmt.Close()

	l := slog.Default().With("query", "UPDATE")
	if len(vals) != 0 {
		l = l.With("values", vals)
	}
	l.Debug(query)

	if _, err = stmt.Exec(vals...); err != nil {
		return errorx.Wrap(err)
	}
	return nil
}
