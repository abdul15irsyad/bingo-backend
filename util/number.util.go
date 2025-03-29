package util

import (
	"math/rand"
)

func RandomInt(min int, max int) int {
	if min > max {
		panic("min cannot be greater than max")
	}

	return rand.Intn(max-min+1) + min
}
