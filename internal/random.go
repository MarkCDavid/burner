package internal

import (
	"math"
	"math/rand"
	"time"
)

type Rng struct {
	_rng *rand.Rand
}

func CreateRandom(seed *int64) *Rng {
	if seed == nil {
		_seed := time.Now().UnixNano()
		seed = &_seed
	}

	return &Rng{
		_rng: rand.New(rand.NewSource(*seed)),
	}
}

func (r *Rng) Expovariate(lambda float64) float64 {
	return -math.Log(1.0-r.float()) / lambda
}

func (r *Rng) float() float64 {
	return 1.0 - r._rng.Float64()
}
