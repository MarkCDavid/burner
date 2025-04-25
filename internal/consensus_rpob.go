package internal

func AddConsensus_RPoB(node *Node) {
	configuration := node.Simulation.Configuration.RazerProofOfBurn

	if !configuration.Enabled {
		return
	}

	node.Consensus = append(node.Consensus, &Consensus_RPoB{
		Enabled: configuration.Enabled,

		Node: node,

		Interval:   configuration.Interval,
		Difficulty: float64(1),
	})
}

type Consensus_RPoB_Configuration struct {
	Enabled  bool    `yaml:"enabled"`
	Interval float64 `yaml:"interval"`
}

type Consensus_RPoB struct {
	Enabled bool

	Node *Node

	Interval   float64
	Difficulty float64
}

func (c *Consensus_RPoB) Initialize() {
	c.Difficulty = c.Interval * float64(len(c.Node.Simulation.Nodes))
}

func (c *Consensus_RPoB) GetType() ConsensusType {
	return ProofOfBurn
}

func (c *Consensus_RPoB) CanMine(receivedEvent *Event) bool {
	if !c.Enabled {
		return false
	}

	chance := float64(1) / c.Difficulty
	if chance > 1 {
		chance = 0.995
	}
	return c.Node.Simulation.Random.Chance(chance)
}

func (c *Consensus_RPoB) GetNextMiningTime(event *Event) float64 {
	// Computing 1 hash takes barely any time.
	return c.Node.Simulation.CurrentTime + 1
}

func (c *Consensus_RPoB) Synchronize(consensus Consensus) {
	_, ok := consensus.(*Consensus_RPoB)
	if !ok {
		panic("not a rezer proof of burn difficulty")
	}
}

func (c *Consensus_RPoB) Adjust(event *Event) {

}
