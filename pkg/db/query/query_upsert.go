package query

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/cresendoo/decidash-backend/pkg/errorx"
)

type UpsertResult int

const (
	UpsertResultNoAffected UpsertResult = iota
	UpsertResultInserted
	UpsertResultUpdated
)

func AddOrSetStruct(db Conn, v Model) (UpsertResult, error) {
	cols := getStructColumns(v)
	baseQuery, cols := makeQueryString(v.TableName(), cols, stmtInsert)
	updateQuery, _ := makeQueryString(v.TableName(), cols, stmtUpsert)

	stmt, err := db.Prepare(baseQuery + updateQuery)
	if err != nil {
		return UpsertResultNoAffected, errorx.Wrap(err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(getStructValues(cols, v)...)
	if err != nil {
		return UpsertResultNoAffected, errorx.Wrap(err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return UpsertResultNoAffected, errorx.Wrap(err)
	}

	return UpsertResult(rowsAffected), nil
}

func AddOrSetStructs(db Conn, v any) error {
	return splitSliceQuery(db, v, addOrSetStructs)
}

func AddOrCustomSetStruct(db Conn, v Model, setQuery *QueryString) (UpsertResult, error) {
	cols := getStructColumns(v)
	baseQuery, cols := makeQueryString(v.TableName(), cols, stmtInsert)
	query := fmt.Sprintf("%s \nON DUPLICATE KEY UPDATE %s", baseQuery, strings.Join(setQuery.OnDuplicates, ","))
	stmt, err := db.Prepare(query)
	if err != nil {
		return UpsertResultNoAffected, errorx.Wrap(err)
	}
	defer stmt.Close()
	vals := getStructValues(cols, v)
	result, err := stmt.Exec(append(vals, setQuery.Vals...)...)
	if err != nil {
		return UpsertResultNoAffected, errorx.Wrap(err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return UpsertResultNoAffected, errorx.Wrap(err)
	}

	return UpsertResult(rowsAffected), nil
}

func addOrSetStructs(db Conn, slice reflect.Value) error {
	baseType := indirectType(slice.Type().Elem())
	baseValue := reflect.New(baseType).Interface().(Model)

	baseCols := getStructColumns(baseValue)
	baseQuery, cols := makeQueryString(baseValue.TableName(), baseCols, stmtInsertMulti)
	updateStmt, _ := makeQueryString(baseValue.TableName(), baseCols, stmtUpsert)

	valueStmt := strings.Join(stringRepeat(makeValuesString(len(cols)), slice.Len()), ", ")

	stmt, err := db.Prepare(baseQuery + valueStmt + updateStmt)
	if err != nil {
		return errorx.Wrap(err)
	}
	defer stmt.Close()

	var values []interface{}
	for i := 0; i < slice.Len(); i++ {
		values = append(values, getStructValues(cols, slice.Index(i).Interface())...)
	}

	_, err = stmt.Exec(values...)
	if err != nil {
		return errorx.Wrap(err)
	}
	return nil
}
