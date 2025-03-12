package randomness

import (
	"math"
	"math/rand"
	"time"
)

type Randomness struct {
	rng *rand.Rand
}

func NewRandomnessWithSeed(seed int64) *Randomness {
	return &Randomness{
		rng: rand.New(rand.NewSource(seed)),
	}
}

func NewRandomness() *Randomness {
	return &Randomness{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r *Randomness) Next(lowBound, highBound float64) float64 {
	return lowBound + r.float()*(highBound-lowBound)
}

func (r *Randomness) Id() uint64 {
	return uint64(r.rng.Int63())
}

func (r *Randomness) Expovariate(lambda float64) float64 {
	return -math.Log(1.0-r.float()) / lambda
}

func (r *Randomness) float() float64 {
	return 1.0 - r.rng.Float64()
}
