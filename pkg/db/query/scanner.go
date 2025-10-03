package query

import (
	"database/sql"
	"reflect"

	"github.com/cresendoo/decidash-backend/pkg/db/encrypt"
	"github.com/cresendoo/decidash-backend/pkg/errorx"
)

func scanStruct(rows *sql.Rows, st any) error {
	columns, err := rows.Columns()
	if err != nil {
		return errorx.Wrap(err)
	}

	val := reflect.Indirect(reflect.ValueOf(st))
	if val.Kind() != reflect.Struct {
		return errorx.New("Invalid Type Error :(not a Struct)")
	}
	structInfo := getStructInfo(val.Type())

	values := make([]interface{}, 0, len(columns))
	var encValues []*encrypt.Column
	for _, col := range columns {
		if ci, ok := structInfo[col]; ok {
			v := val.Field(ci.structIDX)
			if v.CanAddr() {
				v = v.Addr()
			}
			if ci.Encrypted {
				enc := encrypt.NewColumn(st.(encrypt.EncryptID), v.Interface(), ci.cipher)
				values = append(values, enc)
				encValues = append(encValues, enc)
			} else {
				values = append(values, v.Interface())
			}
		}
	}

	if err := rows.Scan(values...); err != nil {
		return errorx.Wrap(err)
	}

	for _, enc := range encValues {
		if err := enc.Assign(); err != nil {
			return errorx.Wrap(err)
		}
	}
	return nil
}
