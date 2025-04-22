package internal

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type Simulation struct {
	Configuration           Configuration
	Nodes                   []Node
	BlockMinedEventQueue    EventQueue
	BlockReceivedEventQueue EventQueue
	Random                  *Rng
	BlockCount              int
	CurrentTime             float64
}

var simulation Simulation

func Simulate(configuration_path string, seed *int64) error {
	configuration := mustLoadConfiguration(configuration_path)

	random := CreateRandom(seed)
	nodes := make([]Node, configuration.NodeCount)

	simulation = Simulation{
		Configuration: configuration,
		Nodes:         nodes,

		BlockMinedEventQueue:    CreateEventQueue(),
		BlockReceivedEventQueue: CreateEventQueue(),

		Random:      random,
		CurrentTime: 0,
	}

	for nodeId := 0; nodeId < configuration.NodeCount; nodeId += 1 {
		ScheduleBlockMinedEvent(nodeId, 0)
	}

	event := 0
	for simulation.BlockMinedEventQueue.Len() > 0 || simulation.BlockReceivedEventQueue.Len() > 0 {
		event += 1
		time.Sleep(time.Millisecond)

		logrus.Infof("")
		logrus.Infof("")
		logrus.Infof("===========================================")
		logrus.Infof("================= %6d ==================", event)
		logrus.Infof("===========================================")
		logrus.Infof("Block Mined Events: %d | Block Received Events: %d | Pre Pop", simulation.BlockMinedEventQueue.Len(), simulation.BlockReceivedEventQueue.Len())
		logrus.Infof("===========================================")
		for nodeIndex, node := range simulation.Nodes {
			currentEvent, _ := json.Marshal(node.CurrentEvent)
			logrus.Infof("%d - %s", nodeIndex, string(currentEvent))
		}
		logrus.Infof("===========================================")
		blockMinedEvent := simulation.BlockMinedEventQueue.Peek()
		blockReceivedEvent := simulation.BlockReceivedEventQueue.Peek()

		if blockMinedEvent == nil && blockReceivedEvent == nil {
			panic("No events left? Not possible.")
		}

		if blockReceivedEvent == nil || blockMinedEvent.DispatchAt < blockReceivedEvent.DispatchAt {
			blockMinedEvent = simulation.BlockMinedEventQueue.Pop()
			simulation.CurrentTime = blockMinedEvent.DispatchAt

			logrus.Infof("Block Mined Events: %d | Block Received Events: %d | Post Pop", simulation.BlockMinedEventQueue.Len(), simulation.BlockReceivedEventQueue.Len())
			HandleBlockMinedEvent(blockMinedEvent)
		} else if blockMinedEvent == nil || blockReceivedEvent.DispatchAt < blockMinedEvent.DispatchAt {
			var i string
			fmt.Scanln(&i)
			blockReceivedEvent = simulation.BlockReceivedEventQueue.Pop()
			simulation.CurrentTime = blockReceivedEvent.DispatchAt

			logrus.Infof("Block Mined Events: %d | Block Received Events: %d | Post Pop", simulation.BlockMinedEventQueue.Len(), simulation.BlockReceivedEventQueue.Len())
			HandleBlockReceivedEvent(blockReceivedEvent)
		}

		if simulation.CurrentTime > simulation.Configuration.SimulationTime {
			break
		}
	}

	return nil
}
