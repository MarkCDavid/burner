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

	BlockCount       int64
	TransactionCount int64

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
		Statistics:  Statistics{},
	}
}

func (s *Simulation) AdvanceTimeTo(time float64) {
	s.CurrentTime = time
	s.ProgressBar.SetCurrent(int64(s.CurrentTime))
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
		event := s.Events.Pop()

		if event == nil {
			panic("no events")
		}

		s.AdvanceTimeTo(event.DispatchAt)

		switch event.EventType {
		case BlockMinedEvent:
			s.HandleBlockMinedEvent(event)
		case BlockReceivedEvent:
			s.HandleBlockReceivedEvent(event)
		}
	}
	logrus.Infof("Average block mining time: %f", s.Statistics.GetAverageBlockMiningTime())
}
