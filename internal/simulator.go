package internal

import (
	"fmt"
	"os"
)

type Simulation struct {
	Configuration Configuration
	Nodes         []Node
	Blocks        []Block
	Blockchain    []int
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
		Blockchain:  []int{0},
		Queue:       CreateEventQueue(),
		Random:      CreateRandom(seed),
		CurrentTime: 0,
	}

	for nodeId := 0; nodeId < configuration.NodeCount; nodeId += 1 {
		scheduleBlockMinedEvent(nodeId, 0)
	}

	processedEvents := 0
	for simulation.Queue.Len() > 0 {
		processedEvents += 1
		event := simulation.Queue.Pop()

		simulation.CurrentTime = event.DispatchAt

		if simulation.CurrentTime > simulation.Configuration.SimulationTime {
			break
		}

		switch event.Type {
		case BlockMinedEvent:
			// fmt.Printf("%d (%f): Mined | By: %p | Previous: %p\n", processedEvents, event.DispatchAt, event.Node, event.Block)
			handleBlockMinedEvent(event)
		case BlockReceivedEvent:
			// fmt.Printf("%d (%f): Received | By: %p | Block: %p\n", processedEvents, event.DispatchAt, event.Node, event.Block)
			handleBlockReceivedEvent(event)
		default:
			panic("Unknown event")
		}
	}
	fmt.Println("Total blocks worked on:", len(simulation.Blocks))
	fmt.Println("Block in blockchain (including branches):", len(simulation.Blockchain))
	ExportBlocksToDot(simulation.Blockchain, "chain.dot")

	return nil
}

func scheduleBlockMinedEvent(minedBy int, previousBlock int) {
	minedAt := simulation.CurrentTime + simulation.Random.Expovariate(1.0/float64(simulation.Configuration.AverageBlockFrequencyInSeconds*simulation.Configuration.NodeCount))
	block := Block{
		Node:          minedBy,
		PreviousBlock: previousBlock,
		MinedAt:       minedAt,
		Depth:         simulation.Blocks[previousBlock].Depth + 1,
	}

	simulation.Blocks = append(simulation.Blocks, block)
	currentBlock := len(simulation.Blocks) - 1
	simulation.Nodes[minedBy].CurrentlyMinedBlock = currentBlock

	simulation.Queue.Push(&Event{
		Type:       BlockMinedEvent,
		Node:       minedBy,
		Block:      currentBlock,
		DispatchAt: minedAt,
	})
}

func scheduleBlockReceivedEvent(receivedBy int, minedBlock int) {
	offset := 0.0
	if simulation.Configuration.AverageNetworkLatencyInSeconds > 0 {
		offset = simulation.Random.Expovariate(1.0 / float64(simulation.Configuration.AverageNetworkLatencyInSeconds))
	}
	simulation.Queue.Push(&Event{

		Type:       BlockReceivedEvent,
		Node:       receivedBy,
		Block:      minedBlock,
		DispatchAt: simulation.CurrentTime + offset,
	})
}

func handleBlockMinedEvent(event *Event) {
	if simulation.Nodes[event.Node].CurrentlyMinedBlock != event.Block {
		return
	}

	simulation.Blockchain = append(simulation.Blockchain, event.Block)

	for currentNode := 0; currentNode < simulation.Configuration.NodeCount; currentNode += 1 {
		if currentNode == event.Node {
			scheduleBlockMinedEvent(currentNode, event.Block)
		} else {
			scheduleBlockReceivedEvent(currentNode, event.Block)
		}
	}
}

func handleBlockReceivedEvent(event *Event) {
	eventBlock := simulation.Blocks[event.Block]
	currentlyMinedBlock := simulation.Blocks[simulation.Nodes[event.Node].CurrentlyMinedBlock]

	if currentlyMinedBlock.PreviousBlock == eventBlock.PreviousBlock {
		scheduleBlockMinedEvent(event.Node, event.Block)
		return
	}

	if currentlyMinedBlock.Depth+simulation.Configuration.ChainReogranizationThreshold <= eventBlock.Depth {
		scheduleBlockMinedEvent(event.Node, event.Block)
		return
	}
}

func ExportBlocksToDot(blockchain []int, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintln(file, "digraph G {")
	// fmt.Fprintln(file, "    rankdir=LR;")

	// Add nodes
	for _, blockId := range blockchain {
		block := simulation.Blocks[blockId]
		blockAddr := fmt.Sprintf("block_%d", blockId)
		nodeAddr := fmt.Sprintf("%d", block.Node)
		label := fmt.Sprintf("Mined by: %s\\nMined at: %f\\nDepth: %d", nodeAddr, block.MinedAt, block.Depth)
		fmt.Fprintf(file, "    \"%s\" [label=\"%s\"];\n", blockAddr, label)
	}

	// Add edges
	for _, blockId := range blockchain {
		block := simulation.Blocks[blockId]
		if block.PreviousBlock != -1 {
			currentAddr := fmt.Sprintf("block_%d", blockId)
			prevAddr := fmt.Sprintf("block_%d", block.PreviousBlock)
			fmt.Fprintf(file, "    \"%s\" -> \"%s\";\n", prevAddr, currentAddr)
		}
	}

	fmt.Fprintln(file, "}")

	return nil
}
