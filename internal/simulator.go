package internal

import "fmt"

type Simulation struct {
	Configuration Configuration
	Nodes         []Node
	Blocks        []Block
	Queue         EventQueue
	Random        *Rng
	CurrentTime   float64
}

var simulation Simulation

func Simulate(configuration_path string, seed *int64) error {
	configuration, err := LoadConfiguration(configuration_path)
	if err != nil {
		return err
	}

	simulation = Simulation{
		Configuration: configuration,
		Nodes:         make([]Node, configuration.NodeCount),
		Blocks: []Block{
			{
				Node:          -1,
				PreviousBlock: -1,
				Depth:         0,
			},
		},
		Queue:       CreateEventQueue(),
		Random:      CreateRandom(seed),
		CurrentTime: 0,
	}

	for nodeId := 0; nodeId < configuration.NodeCount; nodeId += 1 {
		ScheduleBlockMinedEvent(nodeId, 0)
	}

	truAverage := 0
	average := simulation.Configuration.NodeCount * simulation.Configuration.NodeCount

	fmt.Printf("running with %d nodes.\n", simulation.Configuration.NodeCount)
	fmt.Printf("expecting to average %d events, with maximum around %d\n", average, 2*average)
	min := average
	max := average
	processedEvents := 0
	for simulation.Queue.Len() > 0 {

		truAverage += simulation.Queue.Len()
		if processedEvents >= 50000 {
			if simulation.Queue.Len() > max {
				max = simulation.Queue.Len()
			}

			if simulation.Queue.Len() < min {
				min = simulation.Queue.Len()
			}
		}
		processedEvents += 1
		event := simulation.Queue.Pop()

		simulation.CurrentTime = event.DispatchAt

		switch event.Type {
		case BlockMinedEvent:
			// fmt.Printf("%d (%f): Mined | By: %p | Previous: %p\n", processedEvents, event.DispatchAt, event.Node, event.Block)
			HandleBlockMinedEvent(event)
		case BlockReceivedEvent:
			// fmt.Printf("%d (%f): Received | By: %p | Block: %p\n", processedEvents, event.DispatchAt, event.Node, event.Block)
			HandleBlockReceivedEvent(event)
		default:
			panic("Unknown event")
		}

		if simulation.CurrentTime > simulation.Configuration.SimulationTime {
			break
		}
	}
	fmt.Printf("total events processed: %d\n", processedEvents)
	// fmt.Printf("total capacity of blocks is: %d MB\n", (len(simulation.Blocks)*int(BlockSize))/(1024*1024))
	// fmt.Printf("total capacity of events is: %d KB\n", (2*max*int(EventSize))/(1024))
	// fmt.Printf("expecting on average %d events\n", average)
	fmt.Printf("actually on average %f events, min: %d, max: %d\n", float64(truAverage)/float64(processedEvents), min, max)
	// fmt.Printf("random shit, go: minux %f, dividus1: %f, dividus2: %f\n", float64(truAverage)/float64(processedEvents)-float64(average)-float64(simulation.Configuration.NodeCount), float64(min)/float64(average), float64(max)/float64(average))
	// fmt.Println()
	// fmt.Println()
	// fmt.Println()

	// CalculateStatistics(simulation.Blocks)
	ExportBlocksToDotGraph(simulation.Blocks, true, "chain.dot")

	return nil
}
