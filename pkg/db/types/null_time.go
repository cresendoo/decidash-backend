package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"
)

type NullTime struct {
	t sql.NullTime
}

func (t *NullTime) Set(setTime time.Time) {
	t.t.Time = setTime
	t.t.Valid = true
}

func (t *NullTime) Reset() {
	t.t.Time = time.Time{}
	t.t.Valid = false
}

func (t NullTime) IsNull() bool {
	return !t.t.Valid
}

func (t NullTime) String() string {
	return time.Time(t.t.Time).String()
}

func (t NullTime) MarshalJSON() ([]byte, error) {
	if t.t.Valid {
		return json.Marshal(t.t.Time)
	}
	return json.Marshal(nil)
}

func (t *NullTime) Scan(value any) (err error) {
	return t.t.Scan(value)
}

func (t NullTime) Value() (driver.Value, error) {
	return t.t.Value()
}
