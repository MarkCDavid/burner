package internal

func AddConsensus_SPoB(node *Node) {
	configuration := node.Simulation.Configuration.SlimcoinProofOfBurn

	if !configuration.Enabled {
		return
	}

	node.Consensus = append(node.Consensus, &Consensus_SPoB{
		Enabled: configuration.Enabled,

		Node: node,
	})
}

type Consensus_SPoB_Configuration struct {
	Enabled bool `yaml:"enabled"`
}

type Consensus_SPoB struct {
	Enabled bool

	Node *Node
}

func (d *Consensus_SPoB) Initialize() {}

func (d *Consensus_SPoB) GetType() ConsensusType {
	return ProofOfBurn
}

func (d *Consensus_SPoB) CanMine(receivedEvent *Event) bool {
	if !d.Enabled {
		return false
	}

	if receivedEvent.Block.Consensus.GetType() != ProofOfWork {
		return false
	}

	// TODO: Better modeling for success
	// IDEA: Maybe we count ratio between PoB and PoW
	//       if that ratio drops too low, we reduce difficulty,
	//       if it climbs too high, we increse difficulty?
	chance := float64(1.0 / (5.0 * len(d.Node.Simulation.Nodes)))
	return d.Node.Simulation.Random.Chance(chance)
}

func (c *Consensus_SPoB) GetNextMiningTime(event *Event) float64 {
	// Computing 1 hash takes barely any time.
	return c.Node.Simulation.CurrentTime + 0.0001
}

func (c *Consensus_SPoB) Synchronize(consensus Consensus) {
	_, ok := consensus.(*Consensus_SPoB)
	if !ok {
		panic("not a slimcoin proof of burn difficulty")
	}
}

func (c *Consensus_SPoB) Adjust(event *Event) {}
