package matrix

import (
	"math"
)

var (
	// Positive Checks if a float64 pointer is greater than zero.
	Positive = func(i *float64) bool {
		return *i > 0
	}

	// Negative Checks if a float64 pointer is less than zero.
	Negative = func(i *float64) bool {
		return *i < 0
	}

	// Odd Checks if a float64 pointer is not exactly divisible by 2.0.
	Odd = func(i *float64) bool {
		return math.Mod(*i, 2.0) != 0.0
	}

	// Even Checks if a float64 pointer is exactly divisible by zero.
	Even = func(i *float64) bool {
		return math.Mod(*i, 2.0) == 0.0
	}
)
