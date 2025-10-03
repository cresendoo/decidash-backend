package query

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/cresendoo/decidash-backend/pkg/errorx"
)

func AddStruct(db Conn, v Model, options ...addStructOptionFunc) error {
	opt := &addStructOption{}
	for _, f := range options {
		f(opt)
	}

	cols := getStructColumns(
		v,
		excludeAutoIncrement(true),
		withManagedTimeColumn(opt.withManagedTimeColumn),
	)
	baseQuery, cols := makeQueryString(v.TableName(), cols, stmtInsert)

	stmt, err := db.Prepare(baseQuery)
	if err != nil {
		return errorx.Wrap(err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(getStructValues(cols, v)...)
	if err != nil {
		return errorx.Wrap(err)
	}
	if id, err := result.LastInsertId(); err == nil {
		if av, ok := v.(AutoIncrementModel); ok {
			av.SetID(id)
		}
	}
	return nil
}

type addStructOption struct {
	withManagedTimeColumn bool
}

type addStructOptionFunc func(*addStructOption)

func WithTimestamp() addStructOptionFunc {
	return func(o *addStructOption) { o.withManagedTimeColumn = true }
}

// AddStructs : Insert multiple struct
func AddStructs(db Conn, v interface{}) error {
	return splitSliceQuery(db, v, addStructs)
}

func addStructs(db Conn, slice reflect.Value) error {
	baseType := indirectType(slice.Type().Elem())
	baseValue := reflect.New(baseType).Interface().(Model)

	cols := getStructColumns(baseValue, excludeAutoIncrement(true))
	baseQuery, cols := makeQueryString(baseValue.TableName(), cols, stmtInsertMulti)

	valueStmt := strings.Join(stringRepeat(makeValuesString(len(cols)), slice.Len()), ", ")

	stmt, err := db.Prepare(fmt.Sprintf("%s%s", baseQuery, valueStmt))
	if err != nil {
		return errorx.Wrap(err)
	}
	defer stmt.Close()

	var values []interface{}
	for i := 0; i < slice.Len(); i++ {
		values = append(values, getStructValues(cols, slice.Index(i).Interface())...)
	}

	result, err := stmt.Exec(values...)
	if err != nil {
		return errorx.Wrap(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errorx.Wrap(err)
	}
	if rowsAffected != int64(slice.Len()) {
		return errorx.New("Slice length and affected rows is not equal").
			With("slice", slice.Len()).
			With("row", rowsAffected)
	}
	return nil
}
