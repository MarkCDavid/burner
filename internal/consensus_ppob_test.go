package internal

import (
	"math"
	"testing"
)

func newTestPPoB(t *testing.T) *Consensus_PPoB {
	t.Helper()

	simulation := &Simulation{
		Random:   CreateRandom(42),
		Database: NewSQLite(":memory:"),
	}

	node := &Node{
		Id:            0,
		Simulation:    simulation,
		PreviousBlock: &Block{Depth: 0},
	}

	consensus := &Consensus_PPoB{
		Enabled: true,
		Node:    node,

		Price: 100,

		EpochIndex:  0,
		EpochLength: 1008,

		EpochTimeElapsed: 0,
		EpochTimeAverage: 600,

		SettledPeriod: 10,
		WorkingPeriod: 1200,
	}
	node.ProofOfBurn = consensus

	return consensus
}

func newProofOfWorkBlockMined(depth int64) *Event_BlockMined {
	return &Event_BlockMined{
		Block: &Block{
			Depth:     depth,
			Consensus: &Consensus_PoW{},
		},
	}
}

// A fast epoch (blocks come faster than the target) must increase the price.
func TestAdjustPrice_FastEpoch_IncreasesPrice(t *testing.T) {
	c := newTestPPoB(t)

	// 300 s per block against a 600 s target: deviation = 2.
	c.EpochTimeElapsed = float64(c.EpochLength) * 300

	c.AdjustPrice(newProofOfWorkBlockMined(1))

	if c.Price != 200 {
		t.Fatalf("expected price 200, got %f", c.Price)
	}
}

// A slow epoch (blocks come slower than the target) must decrease the price.
func TestAdjustPrice_SlowEpoch_DecreasesPrice(t *testing.T) {
	c := newTestPPoB(t)

	// 1200 s per block against a 600 s target: deviation = 0.5.
	c.EpochTimeElapsed = float64(c.EpochLength) * 1200

	c.AdjustPrice(newProofOfWorkBlockMined(1))

	if c.Price != 50 {
		t.Fatalf("expected price 50, got %f", c.Price)
	}
}

// The deviation must stay inside [0.25, 4].
func TestAdjustPrice_DeviationIsClamped(t *testing.T) {
	c := newTestPPoB(t)
	c.EpochTimeElapsed = 1 // deviation would be ~600000, clamp to 4
	c.AdjustPrice(newProofOfWorkBlockMined(1))
	if c.Price != 400 {
		t.Fatalf("expected price 400, got %f", c.Price)
	}

	c = newTestPPoB(t)
	c.EpochTimeElapsed = float64(c.EpochLength) * 600 * 100 // deviation 0.01, clamp to 0.25
	c.AdjustPrice(newProofOfWorkBlockMined(1))
	if c.Price != 25 {
		t.Fatalf("expected price 25, got %f", c.Price)
	}
}

// A participation drought (a full epoch of non-PoB blocks) must halve the
// price, one time per drought epoch.
func TestAdjust_Drought_HalvesPrice(t *testing.T) {
	c := newTestPPoB(t)

	for depth := int64(1); depth <= c.EpochLength; depth++ {
		c.Adjust(newProofOfWorkBlockMined(depth))
	}

	if c.Price != 100 {
		t.Fatalf("price must not change before the drought threshold, got %f", c.Price)
	}

	c.Adjust(newProofOfWorkBlockMined(c.EpochLength + 1))

	if c.Price != 50 {
		t.Fatalf("expected price 50 after one drought epoch, got %f", c.Price)
	}

	if c.NonProofOfBurnMined != 0 {
		t.Fatalf("expected drought counter reset, got %d", c.NonProofOfBurnMined)
	}

	// A second full drought epoch halves the price again.
	for depth := int64(0); depth <= c.EpochLength; depth++ {
		c.Adjust(newProofOfWorkBlockMined(c.EpochLength + 2 + depth))
	}

	if c.Price != 25 {
		t.Fatalf("expected price 25 after two drought epochs, got %f", c.Price)
	}
}

// HalvePrice must not let the price reach zero.
func TestHalvePrice_StaysPositive(t *testing.T) {
	c := newTestPPoB(t)
	c.Price = math.SmallestNonzeroFloat64

	c.HalvePrice()

	if c.Price <= 0 {
		t.Fatalf("expected a positive price, got %g", c.Price)
	}
}
