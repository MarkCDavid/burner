package internal

import "github.com/sirupsen/logrus"

func AddConsensus_PoW(node *Node) {
	configuration := node.Simulation.Configuration.ProofOfWork

	if !configuration.Enabled {
		return
	}

	if node.ProofOfWork != nil {
		logrus.Fatal("Multiple Proof of Work layers enabled.")
	}

	node.ProofOfWork = &Consensus_PoW{
		Enabled: configuration.Enabled,

		Node: node,

		Power: node.Simulation.Random.LogNormal(AveragePowerUsage_Node_ProofOfWork),

		EpochIndex:       0,
		EpochLength:      configuration.EpochLength,
		BlockFreqency:    configuration.AverageBlockFrequencyInSeconds,
		EpochTimeElapsed: 0,
		Difficulty:       1,
	}
}

type Consensus_PoW_Configuration struct {
	Enabled                        bool    `yaml:"enabled"`
	EpochLength                    int64   `yaml:"epoch_length"`
	AverageBlockFrequencyInSeconds float64 `yaml:"average_block_frequency_in_seconds"`
}

type Consensus_PoW struct {
	Enabled bool

	Node *Node

	Power float64

	EpochLength int64

	EpochIndex       int64
	EpochTimeElapsed float64

	BlockFreqency float64

	Difficulty float64
}

func (c *Consensus_PoW) Initialize() {
	for _, node := range c.Node.Simulation.Nodes {
		c.Difficulty += node.ProofOfWork.GetPower()
	}
	c.Node.Simulation.Database.SaveProofOfWorkConsensus(c, Initialize)
}

func (c *Consensus_PoW) GetType() ConsensusType {
	return ProofOfWork
}

func (c *Consensus_PoW) GetPower() float64 {
	return c.Power
}

func (c *Consensus_PoW) CanMine(event Event) bool {
	return c.Enabled
}

func (c *Consensus_PoW) GetNextMiningTime(event *Event_BlockMined) float64 {
	lambda := c.Node.ProofOfWork.GetPower() / (c.BlockFreqency * c.Difficulty)
	return c.Node.Simulation.CurrentTime + c.Node.Simulation.Random.Expovariate(lambda)
}

func (c *Consensus_PoW) Synchronize(consensus Consensus) {
	from, ok := consensus.(*Consensus_PoW)
	if !ok {
		panic("not a proof of work difficulty")
	}
	c.EpochIndex = from.EpochIndex
	c.EpochLength = from.EpochLength
	c.BlockFreqency = from.BlockFreqency
	c.EpochTimeElapsed = from.EpochTimeElapsed
	c.Difficulty = from.Difficulty

	c.Node.Simulation.Database.SaveProofOfWorkConsensus(c, Synchronize)
}
func (c *Consensus_PoW) Set(difficulty float64) {
	c.Difficulty = difficulty
}

func (c *Consensus_PoW) Adjust(event Event) {
	blockMinedEvent, ok := event.(*Event_BlockMined)
	if !ok {
		return
	}

	if blockMinedEvent.Block.Consensus.GetType() != ProofOfWork {
		return
	}

	c.EpochIndex += 1
	c.EpochTimeElapsed += event.Duration()

	if c.EpochIndex >= c.EpochLength {
		deviation := (c.BlockFreqency * float64(c.EpochLength)) / c.EpochTimeElapsed
		if deviation > 4 {
			deviation = 4
		}
		if deviation < 0.25 {
			deviation = 0.25
		}

		// logrus.Infof("Epoch Time: %f, Average Time: %f, Epoch Index: %d, Adjustment: %f", c.EpochTimeElapsed, c.EpochTimeElapsed/float64(c.EpochIndex), c.EpochIndex, deviation)

		c.Difficulty *= deviation
		c.EpochIndex = 0
		c.EpochTimeElapsed = 0

		c.Node.Simulation.Database.SaveProofOfWorkConsensus(c, Adjust)
	}
}
