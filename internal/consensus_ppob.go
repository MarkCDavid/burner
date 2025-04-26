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

		Power: node.Simulation.Random.LogNormal(AveragePowerUsage_Node_ProofOfBurn),

		ExpectedBlockFrequency: configuration.AverageBlockFrequencyInSeconds,

		BurnParticipationChance: configuration.ParticipationPercentage,

		SettledPeriod: configuration.SettledPeriod,
		WorkingPeriod: configuration.WorkingPeriod,
	}
}
func (c *Consensus_PPoB) GetType() ConsensusType {
	return ProofOfBurn
}

func (c *Consensus_PPoB) GetPower() float64 {
	return c.Power
}

type Consensus_PPoB_Configuration struct {
	Enabled bool `yaml:"enabled"`

	ParticipationPercentage float64 `yaml:"participation_percentage"`

	AverageBlockFrequencyInSeconds float64 `yaml:"average_block_frequency_in_seconds"`

	PriceEpochLength int64 `yaml:"price_epoch_length"`

	SettledPeriod int64 `yaml:"settled_period"`
	WorkingPeriod int64 `yaml:"working_period"`
}

type Consensus_PPoB struct {
	Enabled bool

	Node *Node

	Power float64

	ExpectedBlockFrequency float64

	Price                float64
	PriceAdjustmentIndex int64

	BurnBudget              float64
	BurnParticipationChance float64

	SettledPeriod int64
	WorkingPeriod int64

	BurnTransactions []BurnTransaction
}

type BurnTransaction struct {
	BurnedAt  int64
	BurnedFor float64
	Power     float64
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
	c.Price = c.Node.Simulation.Random.LogNormal(AverageBurnTransaction_Consensus_PriceProofOfBurn)
	c.BurnBudget = c.Node.Simulation.Random.LogNormal(AverageBurnBudget_Consensus_PriceProofOfBurn)
	c.Burn(0)
}

func (c *Consensus_PPoB) CanMine(event Event) bool {
	randomReceivedEvent, ok := event.(*Event_RandomReceived)
	if !ok {
		return false
	}

	if len(c.BurnTransactions) == 0 {
		return false
	}

	for _, burnTransaction := range c.BurnTransactions {
		multiplier := c.BurnDecay(burnTransaction)
		power := burnTransaction.Power * multiplier
		chance := power / float64(len(randomReceivedEvent.Simulation.Nodes))

		if c.Node.Simulation.Random.Chance(chance) {
			return true
		}
	}

	return false
}

func (c *Consensus_PPoB) GetNextMiningTime(event *Event_BlockMined) float64 {
	return c.Node.Simulation.CurrentTime + 1
}

func (c *Consensus_PPoB) Synchronize(consensus Consensus) {
	other, ok := consensus.(*Consensus_PPoB)
	if !ok {
		panic("not a slimcoin proof of burn difficulty")
	}

	c.Price = other.Price
	c.PriceAdjustmentIndex = other.PriceAdjustmentIndex

}

func (c *Consensus_PPoB) Adjust(event Event) {
	c.BurnTransactions = Filter(c.BurnTransactions, func(t BurnTransaction) bool {
		return c.Delta(t) < c.WorkingPeriod
	})

	_, ok := event.(*Event_BlockMined)
	if ok {
		c.PriceAdjustmentIndex++
		if c.PriceAdjustmentIndex >= c.WorkingPeriod {
			c.AdjustPrice()
			c.PriceAdjustmentIndex = 0
		}
		return
	}

	_, ok = event.(*Event_RandomReceived)
	if ok {
		c.Burn(c.Node.PreviousBlock.Depth)
		return
	}
}

func (c *Consensus_PPoB) AdjustPrice() {
	averageBlockFrequency := c.Node.Simulation.Statistics.BlockMiningTime[ProofOfBurn] / float64(c.Node.Simulation.Statistics.BlocksMined[ProofOfBurn])
	deviation := c.ExpectedBlockFrequency / averageBlockFrequency

	if deviation > 2 {
		deviation = 2
	}
	if deviation < 0.5 {
		deviation = 0.5
	}

	c.Price *= deviation
	if c.Price < 0.00001 {
		c.Price = 0.00001
	}
	if c.Price > 10000000000 {
		c.Price = 10000000000
	}
}

func (c *Consensus_PPoB) Burn(depth int64) {
	if c.Node.ProofOfWork != nil {
		participates := c.Node.Simulation.Random.Chance(c.BurnParticipationChance)
		if !participates {
			return
		}
	}

	var currentSpending float64 = 0
	for _, burnTransaction := range c.BurnTransactions {
		currentSpending += burnTransaction.BurnedFor
	}

	desiredPrice := c.Node.Simulation.Random.LogNormal(AverageBurnTransaction_Consensus_PriceProofOfBurn)
	skip := currentSpending+desiredPrice > c.BurnBudget
	if skip {
		return
	}

	bt := BurnTransaction{
		BurnedAt:  depth,
		BurnedFor: desiredPrice,
		Power:     desiredPrice / c.Price,
	}

	c.BurnTransactions = append(c.BurnTransactions, bt)
}
