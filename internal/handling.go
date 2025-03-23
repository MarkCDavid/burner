package internal

func HandleBlockMinedEvent(event *Event) {
	if simulation.Nodes[event.Node].CurrentlyMinedBlock != event.Block {
		return
	}

	simulation.Blocks[event.Block].FinishedAt = event.DispatchAt
	simulation.Blocks[event.Block].Mined = true

	for currentNode := 0; currentNode < len(simulation.Nodes); currentNode += 1 {
		if currentNode == event.Node {
			ScheduleBlockMinedEvent(currentNode, event.Block)
		} else {
			ScheduleBlockReceivedEvent(currentNode, event.Block)
		}
	}
}

func HandleBlockReceivedEvent(event *Event) {
	minedBlock := simulation.Nodes[event.Node].CurrentlyMinedBlock

	if simulation.Blocks[minedBlock].PreviousBlock == simulation.Blocks[event.Block].PreviousBlock {
		simulation.Blocks[minedBlock].FinishedAt = event.DispatchAt
		ScheduleBlockMinedEvent(event.Node, event.Block)
		return
	}

	if simulation.Blocks[minedBlock].Depth+simulation.Configuration.ChainReogranizationThreshold <= simulation.Blocks[event.Block].Depth {
		simulation.Blocks[minedBlock].FinishedAt = event.DispatchAt
		ScheduleBlockMinedEvent(event.Node, event.Block)
		return
	}
}
