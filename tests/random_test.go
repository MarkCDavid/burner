package tests

import (
	"testing"

	pkgRandom "github.com/MarkCDavid/burner/internal/random"
)

func TestNext(t *testing.T) {
	r := pkgRandom.NewRandomnessWithSeed(42)
	low, high := 1.0, 10.0
	val := r.Next(low, high)

	if val < low || val > high {
		t.Errorf("NextDouble() returned value out of bounds: got %f, want between %f and %f", val, low, high)
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
