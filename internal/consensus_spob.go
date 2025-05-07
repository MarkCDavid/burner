package internal

import "github.com/sirupsen/logrus"

func AddConsensus_SPoB(node *Node) {
	configuration := node.Simulation.Configuration.SlimcoinProofOfBurn

	if !configuration.Enabled {
		return
	}

	if node.ProofOfBurn != nil {
		logrus.Fatal("Multiple Proof of Burn layers enabled.")
	}

	node.ProofOfBurn = &Consensus_SPoB{
		Enabled: configuration.Enabled,

		Node: node,

		BlocksMined: *NewSWCounterInt64[ConsensusType](1008),

		Ratio:      float64(1) / configuration.Interval,
		Difficulty: float64(1),
	}
}

type Consensus_SPoB_Configuration struct {
	Enabled  bool    `yaml:"enabled"`
	Interval float64 `yaml:"interval"`
}

type Consensus_SPoB struct {
	Enabled bool

	Node *Node

	Power float64

	BlocksMined SWCounterInt64[ConsensusType]

	Ratio      float64
	Difficulty float64
}

func (c *Consensus_SPoB) GetPower() float64 {
	return c.Power
}
func (c *Consensus_SPoB) GetType() ConsensusType {
	return ProofOfBurn
}

func (c *Consensus_SPoB) Initialize() {
	c.Node.Simulation.Database.SaveSlimcoinProofOfBurnConsensus(c, Initialize)

}

func (c *Consensus_SPoB) CanMine(event Event) bool {
	if !c.Enabled {
		return false
	}

	blockReceivedEvent, ok := event.(*Event_BlockReceived)
	if !ok {
		return false
	}

	if blockReceivedEvent.Block.Consensus.GetType() != ProofOfWork {
		return false
	}

	chance := float64(1) / c.Difficulty
	return c.Node.Simulation.Random.Chance(chance)
}

func (c *Consensus_SPoB) GetNextMiningTime(event *Event_BlockMined) float64 {
	// Computing 1 hash takes barely any time.
	return c.Node.Simulation.CurrentTime + c.Node.Simulation.Random.Float()*0.5
}

func (c *Consensus_SPoB) Synchronize(consensus Consensus) {
	_, ok := consensus.(*Consensus_SPoB)
	if !ok {
		return
	}
}

func (c *Consensus_SPoB) Adjust(event Event) {
	// We calculate the ratio between PoW and PoB blocks
	blockReceivedEvent, ok := event.(*Event_BlockReceived)
	if ok {
		c.BlocksMined.Add(blockReceivedEvent.Block.Consensus.GetType())
		return
	}

	blockMinedEvent, ok := event.(*Event_BlockMined)
	if !ok {
		return
	}

	// If the block mined is a PoB block, we will adjust difficulty.
	if blockMinedEvent.Block.Consensus.GetType() != ProofOfBurn {
		return
	}

	ratio := float64(c.BlocksMined.Get(ProofOfBurn)) / float64(c.BlocksMined.Get(ProofOfWork))
	deviation := ratio / c.Ratio

	if deviation > 4.0 {
		deviation = 4.0
	} else if deviation < 0.25 {
		deviation = 0.25
	}

	c.Difficulty *= deviation

	c.Node.Simulation.Database.SaveSlimcoinProofOfBurnConsensus(c, Adjust)
}
