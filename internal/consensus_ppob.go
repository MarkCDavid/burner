package internal

// import (
// 	"math"
//
// 	"github.com/sirupsen/logrus"
// )
//
// func AddConsensus_PPoB(node *Node) {
// 	configuration := node.Simulation.Configuration.PricingProofOfBurn
//
// 	if !configuration.Enabled {
// 		return
// 	}
//
// 	node.Consensus = append(node.Consensus, &Consensus_PPoB{
// 		Enabled: configuration.Enabled,
//
// 		Node: node,
//
// 		BlockFrequency: configuration.AverageBlockFrequencyInSeconds,
//
// 		DifficultyEpochLength: configuration.DifficultyEpochLength,
// 		PriceEpochLength:      configuration.PriceEpochLength,
// 		BurnEpochLength:       configuration.BurnEpochLength,
//
// 		BurnParticipationChance: configuration.ParticipationPercentage,
//
// 		SettledPeriod: configuration.SettledPeriod,
// 		WorkingPeriod: configuration.WorkingPeriod,
// 	})
// }
//
// type Consensus_PPoB_Configuration struct {
// 	Enabled bool `yaml:"enabled"`
//
// 	ParticipationPercentage float64 `yaml:"participation_percentage"`
//
// 	AverageBlockFrequencyInSeconds float64 `yaml:"average_block_frequency_in_seconds"`
//
// 	DifficultyEpochLength int64 `yaml:"difficulty_epoch_length"`
// 	PriceEpochLength      int64 `yaml:"price_epoch_length"`
// 	BurnEpochLength       int64 `yaml:"burn_epoch_length"`
//
// 	SettledPeriod int64 `yaml:"settled_period"`
// 	WorkingPeriod int64 `yaml:"working_period"`
// }
//
// type Consensus_PPoB struct {
// 	Enabled bool
//
// 	Node *Node
//
// 	BlockFrequency float64
//
// 	Difficulty                 float64
// 	DifficultyEpochIndex       int64
// 	DifficultyEpochLength      int64
// 	DifficultyEpochTimeElapsed float64
//
// 	Price                 float64
// 	PriceEpochIndex       int64
// 	PriceEpochPoBBlocks   int64
// 	PriceEpochLength      int64
// 	PriceEpochTimeElapsed float64
//
// 	BurnBudget              float64
// 	BurnEpochIndex          int64
// 	BurnEpochLength         int64
// 	BurnParticipationChance float64
//
// 	SettledPeriod int64
// 	WorkingPeriod int64
//
// 	BurnTransactions []BurnTransaction
// }
//
// type BurnTransaction struct {
// 	BurnedAt  int64
// 	BurnedFor float64
// 	Power     float64
// }
//
// func (c *Consensus_PPoB) Delta(e *Event, t BurnTransaction) int64 {
// 	return e.Block.Depth - t.BurnedAt
// }
//
// func (c *Consensus_PPoB) BurnDecay(e *Event, t BurnTransaction) float64 {
// 	delta := c.Delta(e, t)
//
// 	if delta < c.SettledPeriod || delta > c.WorkingPeriod {
// 		return 0
// 	}
//
// 	numerator := float64(delta) - float64(c.SettledPeriod)
// 	denominator := float64(c.WorkingPeriod - c.SettledPeriod)
//
// 	return math.Pow(numerator/denominator, 3)
// }
//
// func (c *Consensus_PPoB) Initialize() {
// 	c.Difficulty = 1
// 	c.Price = c.Node.Simulation.Random.LogNormal(AverageBurnTransaction_Consensus_PriceProofOfBurn)
// 	c.BurnBudget = c.Node.Simulation.Random.LogNormal(AverageBurnBudget_Consensus_PriceProofOfBurn)
// }
//
// func (c *Consensus_PPoB) GetType() ConsensusType {
// 	return ProofOfBurn
// }
//
// func (c *Consensus_PPoB) CanMine(receivedEvent *Event) bool {
// 	if !c.Enabled {
// 		return false
// 	}
//
// 	return len(c.BurnTransactions) > 0
// }
//
// func (c *Consensus_PPoB) GetNextMiningTime(event *Event) float64 {
// 	var miningTime float64 = math.Inf(1)
// 	for _, burnTransaction := range c.BurnTransactions {
// 		multiplier := c.BurnDecay(event, burnTransaction)
// 		lambda := burnTransaction.Power / (c.BlockFrequency * c.Difficulty)
// 		delta := c.Node.Simulation.Random.Expovariate(multiplier * lambda)
//
// 		// logrus.Infof("%d burning - multiplier %f, power %f, delta %f", c.Node.Id, multiplier, burnTransaction.Power, delta)
// 		if delta < miningTime {
// 			miningTime = delta
// 		}
// 	}
//
// 	return c.Node.Simulation.CurrentTime + miningTime
// }
//
// func (c *Consensus_PPoB) Synchronize(consensus Consensus) {
// 	_, ok := consensus.(*Consensus_PPoB)
// 	if !ok {
// 		panic("not a slimcoin proof of burn difficulty")
// 	}
// }
//
// func (c *Consensus_PPoB) Adjust(event *Event) {
// 	Filter(c.BurnTransactions, func(t BurnTransaction) bool {
//
// 		return c.BurnDecay(event, t) > 0.1
// 	})
//
// 	if c.GetType() == event.Block.Consensus.GetType() {
// 		c.AdjustDifficulty(event)
// 	}
// 	c.AdjustPrice(event)
// 	c.Burn(event.Block.Depth)
// }
//
// func (c *Consensus_PPoB) AdjustDifficulty(event *Event) {
// 	c.DifficultyEpochIndex++
// 	c.DifficultyEpochTimeElapsed += event.Duration()
//
// 	if c.DifficultyEpochIndex < c.DifficultyEpochLength {
// 		return
// 	}
// 	deviation := (c.BlockFrequency * float64(c.DifficultyEpochLength)) / c.DifficultyEpochTimeElapsed
// 	if deviation > 4 {
// 		deviation = 4
// 	}
// 	if deviation < 0.25 {
// 		deviation = 0.25
// 	}
//
// 	// logrus.Infof("adjusting difficulty (%f) by %f", c.Difficulty, deviation)
//
// 	c.Difficulty *= deviation
// 	if c.Difficulty < 1 {
// 		c.Difficulty = 1
// 	}
// 	c.DifficultyEpochIndex = 0
// 	c.DifficultyEpochTimeElapsed = 0
// }
//
// func (c *Consensus_PPoB) AdjustPrice(event *Event) {
// 	c.PriceEpochIndex++
// 	if event.Block.Consensus.GetType() == ProofOfBurn {
// 		c.PriceEpochTimeElapsed += event.Duration()
// 	}
//
// 	if c.PriceEpochIndex < c.PriceEpochLength {
// 		return
// 	}
//
// 	deviation := c.PriceEpochTimeElapsed / (c.BlockFrequency * float64(c.PriceEpochLength))
// 	if deviation > 4 {
// 		deviation = 4
// 	}
// 	if deviation < 0.25 {
// 		deviation = 0.25
// 	}
//
// 	c.Price *= deviation
// 	if c.Price < 1 {
// 		c.Price = 1
// 	}
// 	c.PriceEpochIndex = 0
// 	c.PriceEpochTimeElapsed = 0
// }
//
// func (c *Consensus_PPoB) Burn(depth int64) {
// 	c.BurnEpochIndex++
//
// 	if c.BurnEpochIndex < c.BurnEpochLength {
// 		return
// 	}
//
// 	c.BurnEpochIndex = 0
//
// 	participates := c.Node.Simulation.Random.Chance(c.BurnParticipationChance)
// 	if !participates {
// 		return
// 	}
//
// 	var currentSpending float64 = 0
// 	for _, burnTransaction := range c.BurnTransactions {
// 		currentSpending += burnTransaction.BurnedFor
// 	}
//
// 	desiredPrice := c.Node.Simulation.Random.LogNormal(AverageBurnTransaction_Consensus_PriceProofOfBurn)
//
// 	// logrus.Infof("%d Participating. Desired price: %f. Spending %f. Can afford: %f.", c.Node.Id, desiredPrice, currentSpending, c.BurnBudget-currentSpending)
// 	if currentSpending+desiredPrice > c.BurnBudget {
// 		return
// 	}
// 	// logrus.Infof("%d Bought.", c.Node.Id)
//
// 	c.BurnTransactions = append(c.BurnTransactions, BurnTransaction{
// 		BurnedAt:  depth,
// 		BurnedFor: desiredPrice,
// 		Power:     desiredPrice / c.Price,
// 	})
//
// 	for _, t := range c.BurnTransactions {
// 		logrus.Info(t)
// 	}
// }
