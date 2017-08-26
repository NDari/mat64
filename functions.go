package matrix

import (
	"math"
)

var (
	// Positivef64 Checks if a float64 pointer is greater than zero.
	Positivef64 = func(i *float64) bool {
		return *i > 0
	}

	// Negativef64 Checks if a float64 pointer is less than zero.
	Negativef64 = func(i *float64) bool {
		return *i < 0
	}

	// Oddf64 Checks if a float64 pointer is not exactly divisible by 2.0.
	Oddf64 = func(i *float64) bool {
		return math.Mod(*i, 2.0) != 0.0
	}

	// Evenf64 Checks if a float64 pointer is exactly divisible by zero.
	Evenf64 = func(i *float64) bool {
		return math.Mod(*i, 2.0) == 0.0
	}
)
