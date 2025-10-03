package query

import (
	"database/sql"
	"log/slog"
	"reflect"

	"github.com/cresendoo/decidash-backend/pkg/errorx"
)

func splitSliceQuery(conn Conn, v any, sliceFunc func(conn Conn, slice reflect.Value) error) (err error) {
	batch := 500
	slice := reflect.Indirect(reflect.ValueOf(v))
	if slice.Kind() != reflect.Slice || !slice.CanSet() {
		return errorx.New("sliceFunc value must be pointer of slice")
	}
	// No rows
	l := slice.Len()
	if l == 0 {
		return nil
	}
	if l > batch {
		// split query
		var tx *sql.Tx
		switch c := conn.(type) {
		case Beginner:
			tx, err = c.Begin()
			if err != nil {
				return errorx.Wrap(err)
			}
			defer func() {
				if err != nil {
					if err := tx.Rollback(); err != nil {
						slog.Default().With("error", err).Error("rollback failed")
					}
				} else {
					if err := tx.Commit(); err != nil {
						slog.Default().With("error", err).Error("commit failed")
					}
				}
			}()
		case *sql.Tx:
			tx = c
		default:
			return errorx.New("conn must implement Beginner interface")
		}
		for st := 0; st < l; st += batch {
			ed := l
			if l > st+batch {
				ed = st + batch
			}
			err = sliceFunc(tx, slice.Slice(st, ed))
			if err != nil {
				return
			}
		}
	} else {
		err = sliceFunc(conn, slice)
	}
	return
}
