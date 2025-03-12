package main

import (
	"fmt"
	"os"
)

const (
	NODE_COUNT                   = 300
	MAX_EVENTS                   = 100000
	BLOCK_FREQUENCY_IN_SECONDS   = 600
	BLOCK_PROPOGATION_IN_SECONDS = 6
	DEPTH_NEEDED_FOR_SWITCH      = 1
)

var Nodes [NODE_COUNT]Node = [NODE_COUNT]Node{}
var Blocks []*Block = []*Block{
	{
		Node:          nil,
		PreviousBlock: nil,
		Depth:         0,
	},
}
var Queue EventQueue = CreateEventQueue()
var Random *Rng = CreateRandom(nil)
var SimulationTime = 0.0

func main() {
	for nodeId := 0; nodeId < NODE_COUNT; nodeId += 1 {
		scheduleBlockMinedEvent(&Nodes[nodeId], Blocks[0])
	}

	processedEvents := 0
	for processedEvents < MAX_EVENTS && Queue.Len() > 0 {
		processedEvents += 1
		event := Queue.Pop()

		SimulationTime = event.DispatchAt

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
	fmt.Println(len(Blocks))
	ExportBlocksToDot(Blocks, "chain.dot")
}

func scheduleBlockMinedEvent(minedBy *Node, previousBlock *Block) {
	minedAt := SimulationTime + Random.Expovariate(1.0/BLOCK_FREQUENCY_IN_SECONDS)
	block := &Block{
		Node:          minedBy,
		PreviousBlock: previousBlock,
		MinedAt:       minedAt,
		Depth:         previousBlock.Depth + 1,
	}

	minedBy.CurrentlyMinedBlock = block

	Queue.Push(&Event{
		Type:       BlockMinedEvent,
		Node:       minedBy,
		Block:      block,
		DispatchAt: minedAt,
	})
}

func scheduleBlockReceivedEvent(receivedBy *Node, minedBlock *Block) {
	Queue.Push(&Event{
		Type:       BlockReceivedEvent,
		Node:       receivedBy,
		Block:      minedBlock,
		DispatchAt: SimulationTime + Random.Expovariate(1.0/BLOCK_PROPOGATION_IN_SECONDS),
	})
}

func handleBlockMinedEvent(event *Event) {
	if event.Node.CurrentlyMinedBlock != event.Block {
		return
	}

	Blocks = append(Blocks, event.Block)

	for nodeId := 0; nodeId < NODE_COUNT; nodeId += 1 {
		node := &Nodes[nodeId]
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

	if event.Node.CurrentlyMinedBlock.Depth+DEPTH_NEEDED_FOR_SWITCH <= event.Block.Depth {
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
