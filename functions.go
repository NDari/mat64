package mat64

import (
	"math"
)

var (
	Positive = func(i *float64) bool {
		return *i > 0
	}

	Negative = func(i *float64) bool {
		return *i < 0
	}

	Odd = func(i *float64) bool {
		return math.Mod(*i, 2.0) != 0.0
	}

	Even = func(i *float64) bool {
		return math.Mod(*i, 2.0) == 0.0
	}

	Square = func(i *float64) {
		*i *= *i
	}
)
