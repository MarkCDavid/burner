package internal

import (
	"fmt"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/sirupsen/logrus"
)

type Simulation struct {
	Configuration Configuration

	Nodes  []*Node
	Events EventQueue

	Random *Rng

	BlockCount       int64
	ForkCount        int64
	TransactionCount int64

	CurrentTime float64
	ProgressBar *pb.ProgressBar

	Database *SQLite
}

func NewSimulation(configuration_path string) *Simulation {
	configuration := mustLoadConfiguration(configuration_path)
	random := CreateRandom(configuration.Seed)

	databasePath := fmt.Sprintf("result/%s (%s).sqlite", configuration.Name, time.Now().Format("2006-01-02 15:04:05"))
	return &Simulation{
		Configuration: configuration,
		Nodes:         make([]*Node, 0),
		Events:        CreateEventQueue(),

		Random: random,

		BlockCount: 0,
		ForkCount:  0,

		CurrentTime: 0,
		ProgressBar: pb.StartNew(int(configuration.SimulationTime)),
		Database:    NewSQLite(databasePath),
	}
}

func (s *Simulation) AdvanceTimeTo(time float64) float64 {
	deltaTime := time - s.CurrentTime
	s.TransactionCount += s.Random.Gaussian(float64(s.Configuration.AverageTransactionsPerSecond)*deltaTime, 0.8)
	s.CurrentTime = time
	s.ProgressBar.SetCurrent(int64(s.CurrentTime))
	return deltaTime
}

func (s *Simulation) InitializeNodes() {
	for nodeIndex := int64(0); nodeIndex < s.Configuration.NodeCount; nodeIndex++ {
		node := NewNode(s)
		s.Nodes = append(s.Nodes, node)
	}

	for _, node := range s.Nodes {
		if node.ProofOfWork != nil {
			node.ProofOfWork.Initialize()
		}

		if node.ProofOfBurn != nil {
			node.ProofOfBurn.Initialize()
		}

		s.Database.SaveNode(node)
	}

}

func (s *Simulation) GetCurrentTransactionCount() int64 {
	return int64(s.CurrentTime) * s.Configuration.AverageTransactionsPerSecond
}

func (s *Simulation) Simulate() {

	logrus.Infof("=================")
	logrus.Infof("Simulation seed: %d", s.Random.GetSeed())
	logrus.Infof("Simulation database: %s", s.Database._path)
	s.Database.SaveLabel(s.Configuration.Name)
	s.InitializeNodes()

	genesisBlock := &Block{
		Id:        0,
		Node:      nil,
		Depth:     0,
		Consensus: &Consensus_Genesis{},
	}
	for _, node := range s.Nodes {
		(&Event_BlockReceived{
			Simulation:    s,
			ReceivedBy:    node,
			Block:         genesisBlock,
			PreviousBlock: nil,
		}).Handle()
	}

	if s.Configuration.PricingProofOfBurn.Enabled {
		s.ScheduleEmitRandomEvent(0, 60)
	}

	for iteration := 0; s.CurrentTime < s.Configuration.SimulationTime; iteration++ {
		if s.Events.Len() == 0 {
			logrus.Error("blockchain stuck - no events available")
			break
		}
		event := s.Events.Pop()
		s.AdvanceTimeTo(event.EventTime())
		event.Handle()
	}

	s.ProgressBar.Finish()
	s.Database.Close()

	logrus.Infof("Simulation ran for: %f", s.CurrentTime)
}
