package internal

import (
	"fmt"
	"os"
)

type Simulation struct {
	Configuration Configuration
	Nodes         []Node
	Blocks        []*Block
	Blockchain    []*Block
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

	genesisBlock := &Block{
		Node:          nil,
		PreviousBlock: nil,
		Depth:         0,
	}

	simulation = Simulation{
		Configuration: configuration,
		Nodes:         make([]Node, configuration.NodeCount),
		Blocks:        []*Block{genesisBlock},
		Blockchain:    []*Block{genesisBlock},
		Queue:         CreateEventQueue(),
		Random:        CreateRandom(seed),
		CurrentTime:   0,
	}

	for nodeId := 0; nodeId < configuration.NodeCount; nodeId += 1 {
		scheduleBlockMinedEvent(&simulation.Nodes[nodeId], simulation.Blockchain[0])
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

func scheduleBlockMinedEvent(minedBy *Node, previousBlock *Block) {
	minedAt := simulation.CurrentTime + simulation.Random.Expovariate(1.0/float64(simulation.Configuration.AverageBlockFrequencyInSeconds))
	block := &Block{
		Node:          minedBy,
		PreviousBlock: previousBlock,
		MinedAt:       minedAt,
		Depth:         previousBlock.Depth + 1,
	}

	simulation.Blocks = append(simulation.Blocks, block)

	minedBy.CurrentlyMinedBlock = block

	simulation.Queue.Push(&Event{
		Type:       BlockMinedEvent,
		Node:       minedBy,
		Block:      block,
		DispatchAt: minedAt,
	})
}

func scheduleBlockReceivedEvent(receivedBy *Node, minedBlock *Block) {
	simulation.Queue.Push(&Event{
		Type:       BlockReceivedEvent,
		Node:       receivedBy,
		Block:      minedBlock,
		DispatchAt: simulation.CurrentTime + simulation.Random.Expovariate(1.0/float64(simulation.Configuration.AverageNetworkLatencyInSeconds*simulation.Configuration.NodeCount)),
	})
}

func handleBlockMinedEvent(event *Event) {
	if event.Node.CurrentlyMinedBlock != event.Block {
		return
	}

	simulation.Blockchain = append(simulation.Blockchain, event.Block)

	for nodeId := 0; nodeId < simulation.Configuration.NodeCount; nodeId += 1 {
		node := &simulation.Nodes[nodeId]
		if node == event.Node {
			scheduleBlockMinedEvent(node, event.Block)
		} else {
			scheduleBlockReceivedEvent(node, event.Block)
		}
	}
}

func handleBlockReceivedEvent(event *Event) {
	if event.Node.CurrentlyMinedBlock.PreviousBlock == event.Block.PreviousBlock {
		scheduleBlockMinedEvent(event.Node, event.Block)
		return
	}

	if event.Node.CurrentlyMinedBlock.Depth+simulation.Configuration.ChainReogranizationThreshold <= event.Block.Depth {
		scheduleBlockMinedEvent(event.Node, event.Block)
		return
	}
}

func ExportBlocksToDot(blocks []*Block, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintln(file, "digraph G {")
	// fmt.Fprintln(file, "    rankdir=LR;")

	// Add nodes
	for _, block := range blocks {
		blockAddr := fmt.Sprintf("block_%p", block)
		nodeAddr := fmt.Sprintf("%p", block.Node)
		label := fmt.Sprintf("Mined by: %s\\nMined at: %f\\nDepth: %d", nodeAddr, block.MinedAt, block.Depth)
		fmt.Fprintf(file, "    \"%s\" [label=\"%s\"];\n", blockAddr, label)
	}

	// Add edges
	for _, block := range blocks {
		if block.PreviousBlock != nil {
			currentAddr := fmt.Sprintf("block_%p", block)
			prevAddr := fmt.Sprintf("block_%p", block.PreviousBlock)
			fmt.Fprintf(file, "    \"%s\" -> \"%s\";\n", prevAddr, currentAddr)
		}
	}

	fmt.Fprintln(file, "}")

	return nil
}
