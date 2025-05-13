package internal

import "github.com/sirupsen/logrus"

func AddConsensus_RPoB(node *Node) {
	configuration := node.Simulation.Configuration.RazerProofOfBurn

	if !configuration.Enabled {
		return
	}

	if node.ProofOfBurn != nil {
		logrus.Fatal("Multiple Proof of Burn layers enabled.")
	}

	node.ProofOfBurn = &Consensus_RPoB{
		Enabled: configuration.Enabled,

		Node: node,

		Interval:   configuration.Interval,
		Difficulty: float64(1),
	}
}

type Consensus_RPoB_Configuration struct {
	Enabled  bool    `yaml:"enabled"`
	Interval float64 `yaml:"interval"`
}

type Consensus_RPoB struct {
	Enabled bool

	Node *Node

	Power float64

	Interval   float64
	Difficulty float64
}

func (c *Consensus_RPoB) GetType() ConsensusType {
	return ProofOfBurn
}

func (c *Consensus_RPoB) Initialize() {
	c.Difficulty = c.Interval * float64(len(c.Node.Simulation.Nodes))
	c.Node.Simulation.Database.SaveRazerProofOfBurnConsensus(c, Initialize)
}

func (c *Consensus_RPoB) CanMine(event Event) bool {
	if !c.Enabled {
		return false
	}

	chance := ClampFloat64(float64(1)/c.Difficulty, 0, 1)
	return c.Node.Simulation.Random.Chance(chance)
}

func (c *Consensus_RPoB) GetNextMiningTime(event *Event_BlockMined) float64 {
	// Computing 1 hash takes barely any time.
	return c.Node.Simulation.CurrentTime + c.Node.Simulation.Random.Float() + 0.5
}

func (c *Consensus_RPoB) Synchronize(consensus Consensus) {}
func (c *Consensus_RPoB) Adjust(event Event)              {}
