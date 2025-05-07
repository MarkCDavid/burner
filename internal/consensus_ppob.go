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

		// Power: node.Simulation.Random.LogNormal(AveragePowerUsage_Node_ProofOfBurn),

		ExpectedBlockFrequency: configuration.AverageBlockFrequencyInSeconds,

		BurnParticipationChance: configuration.ParticipationPercentage,

		WindowTime:    *NewSWFloat64(1024),
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

	Difficulty float64

	BurnBudget              float64
	BurnParticipationChance float64

	SettledPeriod int64
	WorkingPeriod int64

	WindowTime       SWFloat64
	BurnTransactions []BurnTransaction
}

type BurnTransaction struct {
	BurnedBy  *Node
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
	// c.Price = c.Node.Simulation.Random.LogNormal(AverageBurnTransaction_Consensus_PriceProofOfBurn)
	c.Price = 1
	c.BurnBudget = c.Node.Simulation.Random.LogNormal(AverageBurnBudget_Consensus_PriceProofOfBurn)

	c.Difficulty = 10000

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
	c.PriceAdjustmentIndex = other.PriceAdjustmentIndex

	c.Node.Simulation.Database.SavePricingProofOfBurnConsensus(c, Synchronize)
}

func (c *Consensus_PPoB) Adjust(event Event) {
	c.BurnTransactions = Filter(c.BurnTransactions, func(t BurnTransaction) bool {
		return c.Delta(t) < c.WorkingPeriod
	})
	// logrus.Info(len(c.BurnTransactions))

	blockMinedEvent, ok := event.(*Event_BlockMined)
	if ok {
		if blockMinedEvent.Block.Consensus.GetType() != ProofOfBurn {
			return
		}

		c.WindowTime.Add(blockMinedEvent.IntervalDuration())

		// c.AdjustPrice(blockMinedEvent)
		c.PriceAdjustmentIndex++
		// if c.PriceAdjustmentIndex >= c.Node.Simulation.Configuration.PricingProofOfBurn.WorkingPeriod/4 {
		// if c.PriceAdjustmentIndex >= c.Node.Simulation.Configuration.PricingProofOfBurn.WorkingPeriod*8 {
		if c.PriceAdjustmentIndex >= 2016 {
			c.AdjustPrice(blockMinedEvent)
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

func (c *Consensus_PPoB) AdjustPrice(blockMinedEvent *Event_BlockMined) {
	// averageBlockFrequency := c.Node.Simulation.Statistics.BlockIntervalTime[ProofOfBurn] / float64(c.Node.Simulation.Statistics.BlocksMined[ProofOfBurn])

	// averageBlockFrequency := blockMinedEvent.Block.FinishedAt - blockMinedEvent.PreviousBlock.FinishedAt

	deviation := c.ExpectedBlockFrequency / c.WindowTime.Average()
	// deviation := averageBlockFrequency / c.ExpectedBlockFrequency

	// if deviation > 4 {
	// 	deviation = 4
	// }
	// if deviation < 0.25 {
	// 	deviation = 0.25
	// }

	if deviation > 4 {
		deviation = 4
	}
	if deviation < 0.25 {
		deviation = 0.25
	}

	c.Price *= deviation
	c.Price = ClampPositiveFloat64(c.Price)

	// logrus.Info(c.Price)
	c.Node.Simulation.Database.SavePricingProofOfBurnConsensus(c, Adjust)
}

func (c *Consensus_PPoB) Burn(depth int64) {
	// participates := c.Node.Simulation.Random.Chance(c.BurnParticipationChance / float64(10*len(c.Node.Simulation.Nodes)))
	// if len(c.BurnTransactions) > 0 && !participates {
	// 	return
	// }

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
			Power:     1,
		}

		c.BurnTransactions = append(c.BurnTransactions, bt)
		totalSpent += newTotal
		newTotal = totalSpent + c.Price
	}

	if newTotal > c.BurnBudget {
		overBudgetRatio := (newTotal - c.BurnBudget) / c.BurnBudget

		pricePenalty := math.Exp(-overBudgetRatio * 10)

		if !c.Node.Simulation.Random.Chance(pricePenalty) {

			// logrus.Infof("NOT BURN | Spent: %f Willing: %f Budget: %f", totalSpent, willingToSpend, c.BurnBudget)
			return
		}
		// logrus.Infof("BURN | Spent: %f Willing: %f Budget: %f", totalSpent, willingToSpend, c.BurnBudget)

	}

	bt := BurnTransaction{
		BurnedBy:  c.Node,
		BurnedAt:  depth,
		BurnedFor: c.Price,
		Power:     1,
	}

	c.BurnTransactions = append(c.BurnTransactions, bt)
	c.Node.Simulation.Database.SavePricingProofOfBurnBurnTransaction(bt)
}
