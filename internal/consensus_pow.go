package internal

func AddConsensus_PoW(node *Node) {
	configuration := node.Simulation.Configuration.ProofOfWork

	if !configuration.Enabled {
		return
	}

	node.Consensus = append(node.Consensus, &Consensus_PoW{
		Enabled: configuration.Enabled,

		Node: node,

		EpochIndex:       0,
		EpochLength:      configuration.EpochLength,
		BlockFreqency:    configuration.AverageBlockFrequencyInSeconds,
		EpochTimeElapsed: 0,
		Difficulty:       1,
	})
}

type Consensus_PoW_Configuration struct {
	Enabled                        bool    `yaml:"enabled"`
	EpochLength                    int64   `yaml:"epoch_length"`
	AverageBlockFrequencyInSeconds float64 `yaml:"average_block_frequency_in_seconds"`
}

type Consensus_PoW struct {
	Enabled bool

	Node *Node

	EpochLength int64

	EpochIndex       int64
	EpochTimeElapsed float64

	BlockFreqency float64

	Difficulty float64
}

func (c *Consensus_PoW) Initialize() {
	for _, node := range c.Node.Simulation.Nodes {
		c.Difficulty += node.Power[ProofOfWork]
	}
}

func (c *Consensus_PoW) GetType() ConsensusType {
	return ProofOfWork
}

func (c *Consensus_PoW) CanMine(receivedEvent *Event_BlockReceived) bool {
	return c.Enabled
}

func (c *Consensus_PoW) GetNextMiningTime(event *Event_BlockMined) float64 {
	lambda := c.Node.Power[ProofOfWork] / (c.BlockFreqency * c.Difficulty)
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
}
func (c *Consensus_PoW) Set(difficulty float64) {
	c.Difficulty = difficulty
}

func (c *Consensus_PoW) Adjust(event *Event_BlockMined) {
	if event.Block.Consensus.GetType() != c.GetType() {
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

		// logrus.Infof("Epoch Time: %f, Average Time: %f, Epoch Index: %d, Adjustment: %f", c.EpochTimeElapsed, c.EpochTimeElapsed/float64(c.EpochIndex), c.EpochIndex, adjustment)

		c.Difficulty *= deviation
		c.EpochIndex = 0
		c.EpochTimeElapsed = 0
	}
}
