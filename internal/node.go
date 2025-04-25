package internal

import (
	"math"

	"github.com/sirupsen/logrus"
)

func NewNode(s *Simulation) *Node {
	capability := s.Random.LogNormal(AveragePowerUsage_Node_ProofOfWork)
	efficiency := 1 - math.Pow(s.Random.Float(), 4)

	node := &Node{
		Id:         int64(len(s.Nodes)),
		Simulation: s,

		Capability: capability,
		Efficiency: efficiency,
		Power:      capability,

		Consensus: []Consensus{},
	}

	AddConsensus_SPoB(node)
	AddConsensus_PoW(node)

	node.EnsureConsensusLayerCount("Proof of Work", ProofOfWork)
	node.EnsureConsensusLayerCount("Proof of Burn", ProofOfBurn)

	return node
}

func (n *Node) EnsureConsensusLayerCount(name string, consensusType ConsensusType) {
	total := 0
	for _, consensus := range n.Consensus {
		if consensus.GetType() == consensusType {
			total++
		}
	}

	if total > 1 {
		logrus.Fatalf("Too many %s consensus layers enabled (%d).", name, total)
	}
}

func (to *Node) SynchronizeConsensus(from *Node) {
	for consensusIndex := 0; consensusIndex < len(to.Consensus); consensusIndex++ {
		to.Consensus[consensusIndex].Synchronize(from.Consensus[consensusIndex])
	}
}

type Node struct {
	Id int64

	Simulation *Simulation

	Event *Event

	Capability float64
	Efficiency float64
	Power      float64

	Consensus []Consensus

	Transactions int64
}

func (n *Node) GetConsensus(receivedEvent *Event) Consensus {
	for _, consensus := range n.Consensus {
		if consensus.CanMine(receivedEvent) {
			return consensus
		}
	}
	return nil

}
