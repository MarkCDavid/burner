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

func (r *Randomness) Next() float64 {
	return r.rng.Float64()
}

func (r *Randomness) NextDouble(lowBound, highBound float64) float64 {
	return lowBound + r.rng.Float64()*(highBound-lowBound)
}

func (r *Randomness) Binary(probability float64) bool {
	return r.rng.Float64() < probability
}

func (r *Randomness) Expovariate(lambda float64) float64 {
	return -math.Log(1.0-r.rng.Float64()) / lambda
}

func (r *Randomness) NextGaussian() float64 {
	radiusVariable := r.generateRandomFloat()
	angleVariable := r.generateRandomFloat()

	// Box-Muller transform
	return math.Sqrt(-2.0*math.Log(radiusVariable)) * math.Sin(2.0*math.Pi*angleVariable)
}

func (r *Randomness) generateRandomFloat() float64 {
	return 1.0 - r.rng.Float64()
}

