package tests

import (
	"math"
	"testing"

	pkgRandom "github.com/MarkCDavid/burner/internal/random"
)

func TestNewRandomnessWithSeed(t *testing.T) {
	r1 := pkgRandom.NewRandomnessWithSeed(42)
	r2 := pkgRandom.NewRandomnessWithSeed(42)

	if r1.Next() != r2.Next() {
		t.Errorf("Randomness with the same seed should produce the same result")
	}
}

func TestNextDouble(t *testing.T) {
	r := pkgRandom.NewRandomnessWithSeed(42)
	low, high := 1.0, 10.0
	val := r.NextDouble(low, high)

	if val < low || val > high {
		t.Errorf("NextDouble() returned value out of bounds: got %f, want between %f and %f", val, low, high)
	}
}

func TestBinary(t *testing.T) {
	r := pkgRandom.NewRandomnessWithSeed(42)
	prob := 0.5
	count := 0
	n := 1000
	for i := 0; i < n; i++ {
		if r.Binary(prob) {
			count++
		}
	}

	// Expected roughly half the time to be true
	if count < n/2-50 || count > n/2+50 {
		t.Errorf("Binary() did not return expected probability distribution: got %d true values in %d trials", count, n)
	}
}

func TestExpovariate(t *testing.T) {
	r := pkgRandom.NewRandomnessWithSeed(42)
	lambda := 1.5
	val := r.Expovariate(lambda)

	if val < 0 {
		t.Errorf("Expovariate() returned negative value: %f", val)
	}
}

func TestNextGaussian(t *testing.T) {
	r := pkgRandom.NewRandomnessWithSeed(42)
	val := r.NextGaussian()

	if math.IsNaN(val) || math.IsInf(val, 0) {
		t.Errorf("NextGaussian() returned invalid value: %f", val)
	}
}

