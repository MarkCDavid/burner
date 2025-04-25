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
			BlocksMined:           [DifficultyVariants]int64{},
			TransactionsProcessed: [DifficultyVariants]int64{},
			BlockMiningTime:       [DifficultyVariants]float64{},
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
		s.NewNode()
	}

	var totalCapability float64 = 0
	for nodeIndex := 0; nodeIndex < len(s.Nodes); nodeIndex++ {
		totalCapability += s.Nodes[nodeIndex].Capability
	}

	for nodeIndex := 0; nodeIndex < len(s.Nodes); nodeIndex++ {
		for difficultyIndex := 0; difficultyIndex < len(s.Nodes[nodeIndex].Difficulty); difficultyIndex++ {
			s.Nodes[nodeIndex].Difficulty[difficultyIndex].Set(totalCapability)
		}
	}

}

func (s *Simulation) GetCurrentTransactionCount() int64 {
	return int64(s.CurrentTime) * s.Configuration.AverageTransactionsPerSecond
}

func (s *Simulation) Simulate() {
	s.InitializeNodes()

	initialEvent := &Event{
		Node: -1,
		Block: &Block{
			Id:    s.BlockCount,
			Node:  -1,
			Depth: 0,
			Type:  Genesis,
		},
		PreviousBlock: nil,
	}

	for nodeId := int64(0); nodeId < int64(len(s.Nodes)); nodeId += 1 {
		s.ScheduleBlockMinedEvent(nodeId, initialEvent)
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

	logrus.Infof("Total blocks mined: %d", s.Statistics.GetTotalBlocks())
	logrus.Infof("Total PoW blocks mined: %d", s.Statistics.BlocksMined[ProofOfWork])
	logrus.Infof("Total Slimcoin PoB blocks mined: %d", s.Statistics.BlocksMined[SlimcoinProofOfBurn])
	logrus.Info()
	logrus.Infof("Total mining time: %f", s.Statistics.GetTotalMiningTime())
	logrus.Infof("Total PoW mining time: %f", s.Statistics.BlockMiningTime[ProofOfWork])
	logrus.Infof("Total Slimcoin PoB mining time: %f", s.Statistics.BlockMiningTime[SlimcoinProofOfBurn])
	logrus.Info()
	logrus.Infof("Average block mining time: %f", s.Statistics.GetAverageBlockMiningTime())
	logrus.Infof("Average PoW block mining time: %f", s.Statistics.BlockMiningTime[ProofOfWork]/float64(s.Statistics.BlocksMined[ProofOfWork]))
	logrus.Infof("Average Slimcoin PoB block mining time: %f", s.Statistics.BlockMiningTime[SlimcoinProofOfBurn]/float64(s.Statistics.BlocksMined[SlimcoinProofOfBurn]))
	logrus.Info()
	for i := 0; i < len(s.Statistics.PerNode); i++ {
		logrus.Infof("%d node - on average spent %f on mined block.", i, s.Statistics.PerNode[i].TimeOnSuccessfulMining/float64(s.Statistics.PerNode[i].BlocksMined))
	}
	logrus.Info()
	logrus.Infof("Average transactions per block: %f", s.Statistics.GetAverageTransactionsPerBlock())
	logrus.Infof("Simulation total transactions: %d", s.GetCurrentTransactionCount())
	logrus.Infof("Total processed transactions:  %d", s.Statistics.GetTotalTransactions())
	logrus.Info()
	logrus.Infof("Total power used: %f", s.Statistics.GetTotalPowerUsed())
	logrus.Info()
	logrus.Infof("Power used per block: %f", s.Statistics.GetTotalPowerUsed()/float64(s.Statistics.GetTotalBlocks()))
	logrus.Infof("Power used per transaction: %f", s.Statistics.GetTotalPowerUsed()/float64(s.Statistics.GetTotalTransactions()))
	logrus.Infof("Power used per second: %f", s.Statistics.GetTotalPowerUsed()/float64(s.CurrentTime))
}
