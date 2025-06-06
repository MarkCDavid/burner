package internal

import (
	"math"
	"math/rand"
	"time"
)

type Rng struct {
	_rng  *rand.Rand
	_seed int64
}

func CreateRandom(seed int64) *Rng {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}

	return &Rng{
		_rng:  rand.New(rand.NewSource(seed)),
		_seed: seed,
	}
}

func (r *Rng) GetSeed() int64 {
	return r._seed
}

func (r *Rng) Expovariate(lambda float64) float64 {
	return -math.Log(1.0-r.Float()) / lambda
}

func (r *Rng) LogNormal(mean float64) float64 {
	return math.Exp(math.Log(mean) + r.Float() - 0.5)
}

func (r *Rng) Chance(chance float64) bool {
	return r.Float() < chance
}

func (r *Rng) Float() float64 {
	return 1.0 - r._rng.Float64()
}

func (r *Rng) Poisson(expected float64) int64 {
	L := math.Exp(-expected)
	k := int64(0)
	p := 1.0
	for p > L {
		k++
		p *= r._rng.Float64()
	}
	return k - 1
}

func (r *Rng) Gaussian(expected float64, fraction float64) int64 {
	standartDeviation := expected * fraction
	sample := r._rng.NormFloat64()*standartDeviation + expected
	if sample < 0 {
		return 0
	}
	return int64(sample)
}

func (r *Rng) int(min int64, max int64) int64 {
	return min + (r._rng.Int63n(max) - min)
}
