package internal

import "math"

func (s *Simulation) NewNode() {
	capability := s.Random.LogNormal(AveragePowerUsage_Node_ProofOfWork)
	efficiency := 1 - math.Pow(s.Random.Float(), 4)

	node := &Node{
		Capability: capability,
		Efficiency: efficiency,
		Power:      capability,
		Difficulty: [2]Difficulty{},
	}

	node.Difficulty[ProofOfWork] = NewProofOfWorkDifficulty(s.Configuration.ProofOfWork.EpochLength, s.Configuration.ProofOfWork.AverageBlockFrequencyInSeconds)
	node.Difficulty[ProofOfBurn] = NewProofOfWorkDifficulty(s.Configuration.ProofOfWork.EpochLength, s.Configuration.ProofOfWork.AverageBlockFrequencyInSeconds)

	s.Nodes = append(s.Nodes, node)
}

type Node struct {
	CurrentEvent *Event
	Capability   float64
	Efficiency   float64
	Power        float64
	Difficulty   [2]Difficulty
}
