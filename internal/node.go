package internal

import (
	"github.com/sirupsen/logrus"
)

func NewNode(s *Simulation) *Node {

	node := &Node{
		Id:         int64(len(s.Nodes)),
		Simulation: s,

		Consensus: []Consensus{},
	}

	node.Power[ProofOfWork] = s.Random.LogNormal(AveragePowerUsage_Node_ProofOfWork)
	node.Power[ProofOfBurn] = 0
	// s.Random.LogNormal(AveragePowerUsage_Node_ProofOfBurn)

	AddConsensus_RPoB(node)
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

	Power [2]float64

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
