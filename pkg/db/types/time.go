package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/cresendoo/decidash-backend/pkg/errorx"
)

const (
	noLimitTimeValue int64 = -1
	zeroTimeValue    int64 = 0
)

var (
	noLimitTime = Time(time.Date(9998, 12, 31, 23, 59, 59, 0, time.UTC)) // 9998-12-31 23:59:59
	zeroTime    = Time(time.Date(1970, 01, 01, 00, 00, 00, 0, time.UTC)) // 1970-01-01 00:00:00
)

type Time time.Time

func ZeroTime() Time {
	return zeroTime
}

func ComparisonTime(comparison Time) time.Duration {
	rightTime := time.Now()
	return rightTime.Sub(time.Time(comparison))
}

// minus value is comparison time after then now
// plus value is right time after then now
func Comparison(comparison Time) int {
	rightTime := time.Now()
	return int(rightTime.Sub(time.Time(comparison)).Seconds())
}
func ComparisonAB(A, B Time) int {
	return int(time.Time(A).Sub(time.Time(B)).Seconds())
}

func Midnight() Time {
	now := time.Now()
	year, month, day := now.Date()
	return Time(time.Date(year, month, day, 0, 0, 0, 0, time.UTC))
}

func ResetTime() Time {
	now := time.Now()
	year, month, day := now.Date()
	return Time(time.Date(year, month, day+1, 2, 0, 0, 0, time.UTC))
}

func CheckDate(comparison Time, date int) bool {
	now := time.Now()
	year, month, day := now.Date()
	rightTime := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return int(rightTime.Sub(time.Time(comparison)).Seconds()) <= date*24*60*60

}
func (t *Time) Scan(value interface{}) (err error) {
	if value == nil {
		*t = Time{}
		return nil
	}

	if tt, ok := value.(time.Time); ok {
		*t = Time(tt)
		return nil
	}
	return errorx.New("Can't convert to Time").With("value", value)
}

func (t Time) Value() (driver.Value, error) {
	if t.Time().Equal(zeroTime.Time()) {
		return time.Time(zeroTime), nil
	}
	return time.Time(t), nil
}

func (t *Time) UnmarshalJSON(b []byte) error {
	var sec int64
	if err := json.Unmarshal(b, &sec); err != nil {
		return err
	}
	if sec == 0 {
		*t = Time{}
		return nil
	}
	// time이 no limit time과 같을때는 9998-12-31 23:59:59 으로 변환
	switch sec {
	case noLimitTimeValue:
		*t = noLimitTime
		return nil
	case zeroTimeValue:
		*t = noLimitTime
		return nil
	}
	*t = Time(time.Unix(sec, 0))
	return nil
}

func (t Time) String() string {
	return time.Time(t).String()
}

func (t Time) MarshalJSON() ([]byte, error) {
	// time이 no limit time과 같을때는 -1로 Marshal
	switch t.Time() {
	case noLimitTime.Time():
		return json.Marshal(noLimitTimeValue)
	case zeroTime.Time():
		return json.Marshal(zeroTimeValue)
	}
	return json.Marshal(t.ToInt64())
}

func (t Time) Time() time.Time {
	return time.Time(t)
}

func (t Time) ToInt64() int64 {
	tt := time.Time(t)
	if tt.IsZero() {
		return 0
	}
	return tt.Unix()
}
