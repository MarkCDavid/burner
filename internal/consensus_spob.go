package internal

func AddConsensus_SPoB(node *Node) {
	configuration := node.Simulation.Configuration.SlimcoinProofOfBurn

	if !configuration.Enabled {
		return
	}

	node.Consensus = append(node.Consensus, &Consensus_SPoB{
		Enabled: configuration.Enabled,

		Node: node,

		Ratio:      float64(1) / configuration.Interval,
		Difficulty: float64(1),
	})
}

type Consensus_SPoB_Configuration struct {
	Enabled  bool    `yaml:"enabled"`
	Interval float64 `yaml:"interval"`
}

type Consensus_SPoB struct {
	Enabled bool

	Node *Node

	Ratio      float64
	Difficulty float64
}

func (c *Consensus_SPoB) Initialize() {}

func (c *Consensus_SPoB) GetType() ConsensusType {
	return ProofOfBurn
}

func (c *Consensus_SPoB) CanMine(receivedEvent *Event_BlockReceived) bool {
	if !c.Enabled {
		return false
	}

	if receivedEvent.Block.Consensus.GetType() != ProofOfWork {
		return false
	}

	chance := float64(1) / c.Difficulty
	return c.Node.Simulation.Random.Chance(chance)
}

func (c *Consensus_SPoB) GetNextMiningTime(event *Event_BlockMined) float64 {
	// Computing 1 hash takes barely any time.
	return c.Node.Simulation.CurrentTime + 1
}

func (c *Consensus_SPoB) Synchronize(consensus Consensus) {
	_, ok := consensus.(*Consensus_SPoB)
	if !ok {
		panic("not a slimcoin proof of burn difficulty")
	}
}

func (c *Consensus_SPoB) Adjust(event *Event_BlockMined) {
	if event.Block.Consensus.GetType() != c.GetType() {
		return
	}

	powBlocksMined := event.MinedBy.Simulation.Statistics.BlocksMined[ProofOfWork]
	pobBlocksMined := event.MinedBy.Simulation.Statistics.BlocksMined[ProofOfBurn]

	actualRatio := float64(pobBlocksMined) / float64(powBlocksMined)
	deviation := actualRatio / c.Ratio

	if deviation > 4.0 {
		deviation = 4.0
	} else if deviation < 0.25 {
		deviation = 0.25
	}

	c.Difficulty *= deviation
}
