package utils

import (
	"math"
	"math/rand"
	"time"
)

// NowMilli returns current time as an unix millisecond
func NowMilli() int64 {
	return time.Now().UnixNano() / 1000000
}

// ConvertToMilli convert time to unix millisecond
func ConvertToMilli(t time.Time) int64 {
	return t.UnixMilli()
}

// FromMilli convert unix millisecond to Time object
func FromMilli(msec int64) time.Time {
	sec := msec / 1000
	nsec := (msec % 1000) * 1000000
	return time.Unix(sec, nsec)
}

func ExponentialBackoff(tryCount int) time.Duration {
	factor := float64(tryCount)
	ms := time.Duration(math.Exp(factor) + (rand.Float64()*math.Exp(factor+1.0))*1000)
	return ms * time.Millisecond
}

func TimeNowDate(isUTC bool) (int, time.Month, int) {
	now := time.Now()
	if isUTC {
		now = now.UTC()
	}
	return now.Year(), now.Month(), now.Day()
}
