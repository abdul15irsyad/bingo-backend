package util

import (
	"math/rand"
	"time"
)

func RandomDate(start, end time.Time) time.Time {
	min := start.Unix()
	max := end.Unix()
	sec := rand.Int63n(max-min) + min
	return time.Unix(sec, 0)
}
