package genetic

import (
	"math/rand"
)

func Reproduce(male, female []float64, chance float64) []float64 {
	_base := male
	other := female
	if rand.Float64() < 0.5 {
		_base = female
		other = male
	}

	base := make([]float64, len(_base))
	copy(base, _base)

	for i := range base {
		r := rand.Float64()

		if r <= chance {
			base[i] = rand.Float64()
		} else if r < 0.5+chance/2 {
			base[i] = other[i]
		}
	}

	return base
}
