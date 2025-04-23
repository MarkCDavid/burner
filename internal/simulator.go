package internal

import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"runtime"
)

type Simulation struct {
	Configuration Configuration
	Nodes         []Node
	Events        EventQueue
	Random        *Rng
	BlockCount    int
	CurrentTime   float64
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", m.Alloc/1024/1024)
	fmt.Printf("\tTotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
	fmt.Printf("\tSys = %v MiB", m.Sys/1024/1024)
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func Simulate(configuration_path string, seed *int64) error {
	configuration := mustLoadConfiguration(configuration_path)

	random := CreateRandom(seed)
	nodes := make([]Node, configuration.NodeCount)

	s := &Simulation{
		Configuration: configuration,
		Nodes:         nodes,

		Events: CreateEventQueue(),

		Random:      random,
		CurrentTime: 0,
		BlockCount:  1,
	}

	for nodeId := 0; nodeId < configuration.NodeCount; nodeId += 1 {
		s.ScheduleBlockMinedEvent(nodeId, 0, 1)
	}

	bar := pb.StartNew(int(s.Configuration.SimulationTime))
	for iteration := 0; true; iteration++ {
		event := s.Events.Pop()

		if event == nil {
			panic("no events")
		}

		s.CurrentTime = event.DispatchAt

		switch event.EventType {
		case BlockMinedEvent:
			s.HandleBlockMinedEvent(event)
		case BlockReceivedEvent:
			s.HandleBlockReceivedEvent(event)
		}

		if s.CurrentTime > s.Configuration.SimulationTime {
			break
		}

		bar.SetCurrent(int64(s.CurrentTime))
	}

	return nil
}
