package query

import (
	"database/sql"
	"fmt"
	"log/slog"
	"reflect"

	"github.com/cresendoo/decidash-backend/pkg/errorx"
)

type getStructOption struct {
	primaryKey any
	query      *QueryString
	forUpdate  bool
}

type GetStructOptionFunc func(*getStructOption)

func Query(query *QueryString) GetStructOptionFunc {
	return func(o *getStructOption) { o.query = query }
}

func PrimaryKey(key any) GetStructOptionFunc {
	return func(o *getStructOption) { o.primaryKey = key }
}

func ForUpdate() GetStructOptionFunc {
	return func(o *getStructOption) { o.forUpdate = true }
}

func GetStruct(db Conn, v Model, options ...GetStructOptionFunc) error {
	opt := &getStructOption{}
	for _, f := range options {
		f(opt)
	}

	// validate option
	switch {
	case opt.query != nil && opt.primaryKey != nil:
		return errorx.New("Must use either 'query' or 'primary key'")
	}

	cols := getStructColumns(v, withManagedTimeColumn(true), ignoreOmitEmpty(true))
	var query string
	if opt.primaryKey != nil {
		query, _ = makeQueryString(v.TableName(), cols, stmtSelectByID)
	} else {
		query, _ = makeQueryString(v.TableName(), cols, stmtSelect)
	}

	if opt.query != nil {
		query = fmt.Sprintf("%s %s %s", opt.query.PreStmt, query, opt.query.Stmt)
	}
	// apply option
	if opt.forUpdate {
		if _, ok := db.(*sql.Tx); !ok {
			return errorx.New("SELECT FOR UPDATE have to use tx")
		}
		query += " FOR UPDATE"
	}

	// l := slog.Default().With("query", "SELECT")
	// if opt.query != nil && len(opt.query.Vals) != 0 {
	// 	l = l.With("values", opt.query.Vals)
	// }
	// l.Debug(query)

	var args []any
	if opt.primaryKey != nil {
		args = append(args, opt.primaryKey)
	}
	if opt.query != nil && len(opt.query.Vals) != 0 {
		args = append(args, opt.query.Vals...)
	}

	stmt, err := db.Prepare(query)
	if err != nil {
		return errorx.Wrap(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		return errorx.Wrap(err)
	}
	defer rows.Close()
	if !rows.Next() {
		return errorx.Wrap(sql.ErrNoRows)
	}
	if err = scanStruct(rows, v); err != nil {
		return errorx.Wrap(err)
	}
	return nil
}

func GetStructs(db Conn, v any, where *QueryString, options ...GetStructsOptionFunc) error {
	opt := &getStructsOption{}
	for _, f := range options {
		f(opt)
	}

	slice := reflect.Indirect(reflect.ValueOf(v))
	if slice.Kind() != reflect.Slice || !slice.CanSet() {
		return errorx.New("GetStruct's value must be pointer of slice")
	}

	baseType := indirectType(slice.Type().Elem())
	baseValue := reflect.New(baseType).Interface().(Model)
	var cols []string
	if len(opt.columns) != 0 {
		cols = opt.columns
	} else {
		cols = getStructColumns(baseType, withManagedTimeColumn(true), ignoreOmitEmpty(true))
	}

	baseQuery, _ := makeQueryString(baseValue.TableName(), cols, stmtSelect)
	if where == nil {
		where = &QueryString{}
	}
	query := fmt.Sprintf("%s %s %s", where.PreStmt, baseQuery, where.Stmt)
	l := slog.Default().With("query", "SELECT")
	if len(where.Vals) != 0 {
		l = l.With("values", where.Vals)
	}
	l.Debug(query)

	stmt, err := db.Prepare(query)
	if err != nil {
		return errorx.Wrap(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(where.Vals...)
	if err != nil {
		return errorx.Wrap(err)
	}
	defer rows.Close()
	for rows.Next() {
		val := reflect.New(baseType)
		if err := scanStruct(rows, val.Interface()); err != nil {
			return errorx.Wrap(err)
		}
		slice.Set(reflect.Append(slice, val))
	}
	if slice.Len() == 0 {
		return errorx.Wrap(sql.ErrNoRows)
	}
	return nil
}

type getStructsOption struct {
	columns []string
}

type GetStructsOptionFunc func(*getStructsOption)

func WithColumns(columns ...string) GetStructsOptionFunc {
	return func(o *getStructsOption) { o.columns = columns }
}
