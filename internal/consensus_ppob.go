package internal

import (
	"math"

	"github.com/sirupsen/logrus"
)

func AddConsensus_PPoB(node *Node) {
	configuration := node.Simulation.Configuration.PricingProofOfBurn

	if !configuration.Enabled {
		return
	}

	if node.ProofOfBurn != nil {
		logrus.Fatal("Multiple Proof of Burn layers enabled.")
	}

	node.ProofOfBurn = &Consensus_PPoB{
		Enabled: configuration.Enabled,

		Node: node,

		Price:      1,
		BurnBudget: node.Simulation.Random.LogNormal(configuration.AveragePrice * 2),

		EpochIndex:  0,
		EpochLength: configuration.EpochLength,

		EpochTimeElapsed: 0,
		EpochTimeAverage: configuration.AverageBlockFrequencyInSeconds,

		Difficulty: float64(node.Simulation.Configuration.NodeCount) * (configuration.AverageBlockFrequencyInSeconds),

		SettledPeriod: configuration.SettledPeriod,
		WorkingPeriod: configuration.WorkingPeriod,
	}
}
func (c *Consensus_PPoB) GetType() ConsensusType {
	return ProofOfBurn
}

type Consensus_PPoB_Configuration struct {
	Enabled bool `yaml:"enabled"`

	EpochLength int64 `yaml:"epoch_length"`

	AveragePrice                   float64 `yaml:"average_price"`
	AverageBlockFrequencyInSeconds float64 `yaml:"average_block_frequency_in_seconds"`

	SettledPeriod int64 `yaml:"settled_period"`
	WorkingPeriod int64 `yaml:"working_period"`
}

type Consensus_PPoB struct {
	Enabled bool

	Node *Node

	EpochIndex  int64
	EpochLength int64

	EpochTimeElapsed float64
	EpochTimeAverage float64

	NonProofOfBurnMined int64

	Price      float64
	BurnBudget float64

	Difficulty float64

	SettledPeriod int64
	WorkingPeriod int64

	BurnTransactions []BurnTransaction
}

type BurnTransaction struct {
	BurnedBy  *Node
	BurnedAt  int64
	BurnedFor float64
}

func (c *Consensus_PPoB) Delta(t BurnTransaction) int64 {
	if c.Node.PreviousBlock.Depth <= c.SettledPeriod {
		return c.SettledPeriod + 1
	}

	return c.Node.PreviousBlock.Depth - t.BurnedAt
}

func (c *Consensus_PPoB) BurnDecay(t BurnTransaction) float64 {
	delta := c.Delta(t)

	if delta < c.SettledPeriod || delta > c.WorkingPeriod {
		return 0
	}

	numerator := float64(delta) - float64(c.SettledPeriod)
	denominator := float64(c.WorkingPeriod - c.SettledPeriod)

	return 1 - math.Pow(numerator/denominator, 3)
}

func (c *Consensus_PPoB) Initialize() {
	c.Node.Simulation.Database.SavePricingProofOfBurnConsensus(c, Initialize)
	c.Burn(0)
}

func (c *Consensus_PPoB) CanMine(event Event) bool {
	_, ok := event.(*Event_RandomReceived)
	if !ok {
		return false
	}

	if len(c.BurnTransactions) == 0 {
		return false
	}

	for _, burnTransaction := range c.BurnTransactions {
		chance := c.BurnDecay(burnTransaction) / c.Difficulty
		if c.Node.Simulation.Random.Chance(chance) {
			return true
		}
	}

	return false
}

func (c *Consensus_PPoB) GetNextMiningTime(event *Event_BlockMined) float64 {
	return c.Node.Simulation.CurrentTime + c.Node.Simulation.Random.Float() + 0.5
}

func (c *Consensus_PPoB) Synchronize(consensus Consensus) {
	other, ok := consensus.(*Consensus_PPoB)
	if !ok {
		return
	}

	c.Price = other.Price
	c.EpochIndex = other.EpochIndex
	c.EpochTimeElapsed = other.EpochTimeElapsed

	c.Node.Simulation.Database.SavePricingProofOfBurnConsensus(c, Synchronize)
}

func (c *Consensus_PPoB) Adjust(event Event) {
	c.BurnTransactions = Filter(c.BurnTransactions, func(t BurnTransaction) bool {
		return c.Delta(t) < c.WorkingPeriod-c.Node.Simulation.Random.int(0, c.WorkingPeriod/2)
	})

	blockMinedEvent, ok := event.(*Event_BlockMined)
	if ok {

		if blockMinedEvent.Block.Consensus.GetType() != ProofOfBurn {
			if blockMinedEvent.Block.Depth < 20000 {
				return
			}

			c.NonProofOfBurnMined++

			if c.NonProofOfBurnMined > c.EpochLength {
				if c.EpochIndex == 0 {
					c.EpochIndex = 1
					c.EpochTimeElapsed = 4 * c.EpochTimeAverage
				}
				c.AdjustPrice(blockMinedEvent)

				c.EpochIndex = 0
				c.EpochTimeElapsed = 0
			}

			return
		}
		c.NonProofOfBurnMined = 0
		c.EpochIndex++
		c.EpochTimeElapsed += blockMinedEvent.IntervalDuration()

		if c.EpochIndex >= c.EpochLength {
			c.AdjustPrice(blockMinedEvent)

			c.EpochIndex = 0
			c.EpochTimeElapsed = 0
		}
		return
	}

	_, ok = event.(*Event_RandomReceived)
	if ok {
		c.Burn(c.Node.PreviousBlock.Depth)
		return
	}
}

func (c *Consensus_PPoB) AdjustPrice(blockMinedEvent *Event_BlockMined) {
	average := c.EpochTimeElapsed / float64(c.EpochLength)
	deviation := c.EpochTimeAverage / average

	if deviation > 4 {
		deviation = 4
	}
	if deviation < 0.25 {
		deviation = 0.25
	}

	c.Price *= deviation
	c.Price = ClampPositiveFloat64(c.Price)
	c.Node.Simulation.Database.SavePricingProofOfBurnConsensus(c, Adjust)
}

func (c *Consensus_PPoB) Burn(depth int64) {
	var totalSpent float64 = 0
	for _, burnTransaction := range c.BurnTransactions {
		totalSpent += burnTransaction.BurnedFor
	}

	newTotal := totalSpent + c.Price
	for newTotal < c.BurnBudget {
		bt := BurnTransaction{
			BurnedBy:  c.Node,
			BurnedAt:  depth,
			BurnedFor: c.Price,
		}

		c.BurnTransactions = append(c.BurnTransactions, bt)
		totalSpent += newTotal
		newTotal = totalSpent + c.Price
	}

	if newTotal > c.BurnBudget {
		overBudgetRatio := (newTotal - c.BurnBudget) / c.BurnBudget

		pricePenalty := math.Exp(-overBudgetRatio * 100)

		if !c.Node.Simulation.Random.Chance(pricePenalty) {
			return
		}
	}

	bt := BurnTransaction{
		BurnedBy:  c.Node,
		BurnedAt:  depth,
		BurnedFor: c.Price,
	}

	c.BurnTransactions = append(c.BurnTransactions, bt)
	c.Node.Simulation.Database.SavePricingProofOfBurnBurnTransaction(bt)
}
