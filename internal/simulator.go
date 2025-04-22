package internal

import (
	"time"

	"github.com/sirupsen/logrus"
)

// import (
//
//	"encoding/json"
//	"fmt"
//	"time"
//
//	"github.com/sirupsen/logrus"
//
// )
type Simulation struct {
	Configuration Configuration
	Nodes         []Node
	Forks         map[int]EventQueue
	Network       EventQueue
	Random        *Rng
	BlockCount    int
	ForkCount     int
	CurrentTime   float64
}

func (s *Simulation) GetNextEvent() (*Event, EventType) {
	event := s.Network.Peek()
	forkIndex := -1

	for i, queue := range s.Forks {
		if queue.Len() == 0 {
			continue
		}

		candidate := queue.Peek()

		if event == nil || candidate.DispatchAt < event.DispatchAt {
			event = candidate
			forkIndex = i
		}
	}

	if event == nil {
		return nil, -1
	}

	if forkIndex == -1 {
		return s.Network.Pop(), BlockReceivedEvent
	}

	return s.Forks[forkIndex].Pop(), BlockMinedEvent
}

func Simulate(configuration_path string, seed *int64) error {
	configuration := mustLoadConfiguration(configuration_path)

	random := CreateRandom(seed)
	nodes := make([]Node, configuration.NodeCount)

	s := &Simulation{
		Configuration: configuration,
		Nodes:         nodes,

		Forks:   make(map[int]EventQueue),
		Network: CreateEventQueue(),

		Random:      random,
		CurrentTime: 0,
	}

	s.Forks[0] = CreateEventQueue()

	for nodeId := 0; nodeId < configuration.NodeCount; nodeId += 1 {
		s.ScheduleBlockMinedEvent(nodeId, 0, 0)
	}

	iteration := 0
	for {
		iteration += 1

		logrus.Infof("ITERATION | %d", iteration)
		logrus.Infof("FORKS | %d", len(s.Forks))

		for fork := range s.Forks {
			logrus.Infof("FORK | %d | Size: %d", fork, s.Forks[fork].Len())
			if s.Forks[fork].Len() == 0 {
				logrus.Infof("FORK | %d | Empty - deleting", fork)
				delete(s.Forks, fork)
			}
		}
		event, eventType := s.GetNextEvent()
		s.CurrentTime = event.DispatchAt

		time.Sleep(1)

		if event == nil {
			panic("no events")
		}

		switch eventType {
		case BlockMinedEvent:
			logrus.Debugf("HANDLING | Block Mined Event %s", event.ToString())
			s.HandleBlockMinedEvent(event)
		case BlockReceivedEvent:
			logrus.Debugf("HANDLING | Block Received Event %s", event.ToString())
			s.HandleBlockReceivedEvent(event)
		}

		if s.CurrentTime > s.Configuration.SimulationTime {
			break
		}
	}

	// for simulation.BlockMinedEventQueue.Len() > 0 || simulation.BlockReceivedEventQueue.Len() > 0 {
	// 	event += 1
	// 	time.Sleep(time.Millisecond)
	//
	// 	logrus.Infof("")
	// 	logrus.Infof("")
	// 	logrus.Infof("===========================================")
	// 	logrus.Infof("================= %6d ==================", event)
	// 	logrus.Infof("===========================================")
	// 	logrus.Infof("Block Mined Events: %d | Block Received Events: %d | Pre Pop", simulation.BlockMinedEventQueue.Len(), simulation.BlockReceivedEventQueue.Len())
	// 	logrus.Infof("===========================================")
	// 	for nodeIndex, node := range simulation.Nodes {
	// 		currentEvent, _ := json.Marshal(node.CurrentEvent)
	// 		logrus.Infof("%d - %s", nodeIndex, string(currentEvent))
	// 	}
	// 	logrus.Infof("===========================================")
	// 	blockMinedEvent := simulation.BlockMinedEventQueue.Peek()
	// 	blockReceivedEvent := simulation.BlockReceivedEventQueue.Peek()
	//
	// 	if blockMinedEvent == nil && blockReceivedEvent == nil {
	// 		panic("No events left? Not possible.")
	// 	}
	//
	// 	if blockReceivedEvent == nil || blockMinedEvent.DispatchAt < blockReceivedEvent.DispatchAt {
	// 		blockMinedEvent = simulation.BlockMinedEventQueue.Pop()
	// 		simulation.CurrentTime = blockMinedEvent.DispatchAt
	//
	// 		logrus.Infof("Block Mined Events: %d | Block Received Events: %d | Post Pop", simulation.BlockMinedEventQueue.Len(), simulation.BlockReceivedEventQueue.Len())
	// 		HandleBlockMinedEvent(blockMinedEvent)
	// 	} else if blockMinedEvent == nil || blockReceivedEvent.DispatchAt < blockMinedEvent.DispatchAt {
	// 		var i string
	// 		fmt.Scanln(&i)
	// 		blockReceivedEvent = simulation.BlockReceivedEventQueue.Pop()
	// 		simulation.CurrentTime = blockReceivedEvent.DispatchAt
	//
	// 		logrus.Infof("Block Mined Events: %d | Block Received Events: %d | Post Pop", simulation.BlockMinedEventQueue.Len(), simulation.BlockReceivedEventQueue.Len())
	// 		HandleBlockReceivedEvent(blockReceivedEvent)
	// 	}
	//
	// 	if simulation.CurrentTime > simulation.Configuration.SimulationTime {
	// 		break
	// 	}
	// }

	return nil
}
