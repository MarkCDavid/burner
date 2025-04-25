package internal

import (
	"github.com/cheggaaa/pb/v3"
	"github.com/sirupsen/logrus"
)

type Simulation struct {
	Configuration Configuration

	Nodes  []*Node
	Events EventQueue

	Random *Rng

	BlockCount int64

	CurrentTime float64
	ProgressBar *pb.ProgressBar

	Statistics Statistics
}

func NewSimulation(configuration_path string) *Simulation {
	configuration := mustLoadConfiguration(configuration_path)
	random := CreateRandom(configuration.Seed)

	return &Simulation{
		Configuration: configuration,
		Nodes:         make([]*Node, 0),
		Events:        CreateEventQueue(),

		Random: random,

		BlockCount: 1,

		CurrentTime: 0,
		ProgressBar: pb.StartNew(int(configuration.SimulationTime)),
		Statistics: Statistics{
			BlocksMined:           [2]int64{},
			TransactionsProcessed: [2]int64{},
			BlockMiningTime:       [2]float64{},
			PerNode:               make([]NodeStatistics, 0),
		},
	}
}

func (s *Simulation) AdvanceTimeTo(time float64) float64 {
	deltaTime := time - s.CurrentTime
	s.CurrentTime = time
	s.ProgressBar.SetCurrent(int64(s.CurrentTime))
	return deltaTime
}

func (s *Simulation) InitializeNodes() {
	for nodeIndex := int64(0); nodeIndex < s.Configuration.NodeCount; nodeIndex++ {
		node := NewNode(s)
		s.Nodes = append(s.Nodes, node)
		s.Statistics.PerNode = append(s.Statistics.PerNode, NodeStatistics{})
	}

	for _, node := range s.Nodes {
		for _, consensus := range node.Consensus {
			consensus.Initialize()
		}
	}
}

func (s *Simulation) GetCurrentTransactionCount() int64 {
	return int64(s.CurrentTime) * s.Configuration.AverageTransactionsPerSecond
}

func (s *Simulation) Simulate() {

	logrus.Infof("=================")
	logrus.Infof("Simulation seed: %d", s.Random.GetSeed())
	s.InitializeNodes()

	initialEvent := &Event{
		Node: nil,
		Block: &Block{
			Id:        0,
			Node:      nil,
			Depth:     0,
			Consensus: &Consensus_Genesis{},
		},
		PreviousBlock: nil,
	}

	for _, node := range s.Nodes {
		s.ScheduleBlockMinedEvent(node, initialEvent)
	}

	for iteration := 0; s.CurrentTime < s.Configuration.SimulationTime; iteration++ {
		if s.Events.Len() == 0 {
			logrus.Fatal("blockchain stuck - no events available")
		}
		event := s.Events.Pop()

		s.AdvanceTimeTo(event.DispatchAt)

		switch event.EventType {
		case BlockMinedEvent:
			s.HandleBlockMinedEvent(event)
		case BlockReceivedEvent:
			s.HandleBlockReceivedEvent(event)
		}
	}

	s.ProgressBar.Finish()

	logrus.Infof("Simulation ran for: %f", s.CurrentTime)

	logrus.Infof("Total blocks mined: %d", s.Statistics.GetTotalBlocks())
	logrus.Infof("Total PoW blocks mined: %d", s.Statistics.BlocksMined[ProofOfWork])
	logrus.Infof("Total Slimcoin PoB blocks mined: %d", s.Statistics.BlocksMined[ProofOfBurn])
	logrus.Infof("Ratio blocks mined: %f", float64(s.Statistics.BlocksMined[ProofOfBurn])/float64(s.Statistics.BlocksMined[ProofOfWork]))
	logrus.Info()
	logrus.Infof("Total mining time: %f", s.Statistics.GetTotalMiningTime())
	logrus.Infof("Total PoW mining time: %f", s.Statistics.BlockMiningTime[ProofOfWork])
	logrus.Infof("Total Slimcoin PoB mining time: %f", s.Statistics.BlockMiningTime[ProofOfBurn])
	logrus.Info()
	logrus.Infof("Average block mining time: %f", s.Statistics.GetAverageBlockMiningTime())
	logrus.Infof("Average PoW block mining time: %f", s.Statistics.BlockMiningTime[ProofOfWork]/float64(s.Statistics.BlocksMined[ProofOfWork]))
	logrus.Infof("Average Slimcoin PoB block mining time: %f", s.Statistics.BlockMiningTime[ProofOfBurn]/float64(s.Statistics.BlocksMined[ProofOfBurn]))
	logrus.Info()
	for i := 0; i < len(s.Statistics.PerNode); i++ {
		logrus.Infof("%d node - on average spent %f on mined block.", i, s.Statistics.PerNode[i].TimeOnSuccessfulMining/float64(s.Statistics.PerNode[i].BlocksMined))
	}
	logrus.Info()
	logrus.Infof("Average transactions per block: %f", s.Statistics.GetAverageTransactionsPerBlock())
	logrus.Infof("Simulation total transactions: %d", s.GetCurrentTransactionCount())
	logrus.Infof("Total processed transactions:  %d", s.Statistics.GetTotalTransactions())
	logrus.Infof("Transactions per second:  %f", float64(s.Statistics.GetTotalTransactions())/s.CurrentTime)
	logrus.Info()
	logrus.Infof("Total power used: %f", s.Statistics.GetTotalPowerUsed())
	logrus.Info()
	logrus.Infof("Power used per block: %f", s.Statistics.GetTotalPowerUsed()/float64(s.Statistics.GetTotalBlocks()))
	logrus.Infof("Power used per transaction: %f", s.Statistics.GetTotalPowerUsed()/float64(s.Statistics.GetTotalTransactions()))
	logrus.Infof("Power used per second: %f", s.Statistics.GetTotalPowerUsed()/s.CurrentTime)
}
