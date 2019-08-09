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
			base[i] += (rand.Float64() - 0.5) //(base[i] + base[i] + rand.Float64()) / 3
			if base[i] > 1 {
				base[i] = 1
			}
			if base[i] < -1 {
				base[i] = -1
			}

		} else if r < 0.5+chance/2 {
			//if r < 0.25 {
			// base[i] = (base[i] + other[i]) / 2
			// 	continue
			// }
			base[i] = other[i]
			//base[i] = (base[i] + other[i]) / 2
		}
	}

	return base
}
