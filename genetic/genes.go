package genetic

import (
	"math/rand"
	"sort"
)

func Reproduce(parents [][]float64, recombChance, mutChance float64) []float64 {
	lp := len(parents)
	baseChance := rand.Intn(lp)
	_base := parents[baseChance]
	parents[baseChance] = parents[lp-1]
	lp--

	base := make([]float64, len(_base))
	copy(base, _base)

	for i := 0; i < lp; i++ {
		offset := rand.Intn(len(base))
		for n := offset; n < len(base); n++ {
			base[n] = parents[i][n]
		}
	}

	for i := range base {
		r := rand.Float64()
		if r <= mutChance {
			base[i] = rand.Float64() + rand.Float64() - 1
			continue
		}
	}

	return base
}

type Entity interface {
	Fitness() float64
}

func Select(population []Entity) Entity {
	if len(population) == 0 {
		return nil
	}

	min := 0.0
	for i := range population {
		f := population[i].Fitness()
		if f < min {
			min = f
		}
	}

	total := 0.0
	for i := range population {
		total += population[i].Fitness() - min
	}

	sort.Slice(population, func(i, j int) bool {
		return population[i].Fitness() > population[j].Fitness()
	})

	sel := func(i int) Entity {
		e := population[i]
		population = append(population[:i], population[i+1:]...)
		return e
	}

	r := rand.Float64()
	f := 0.0
	for i := range population {
		f += (population[i].Fitness() - min) / total
		if f > r {
			return sel(i)
		}
	}

	return sel(0)
}
