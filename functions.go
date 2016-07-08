package mat64

import (
	"math"
)

const (
	Positive = func(i *float64) bool {
		return &i > 0
	}

	Negative = func(i *float64) bool {
		return &i < 0
	}

	Odd = func(i *float64) bool {
		return math.Mod(&i, 2.0) != 0.0
	}

	Even = func(i *float64) bool {
		return math.Mod(&i, 2.0) == 0.0
	}

	Square = func(i *float64) float64 {
		return &i * &i
	}
)
