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

func HandleBlockMinedEventPoB(event *Event) {
	if pobSimulation.Nodes[event.Node].CurrentlyMinedBlock != event.Block {
		return
	}

	pobSimulation.Blocks[event.Block].FinishedAt = event.DispatchAt
	pobSimulation.Blocks[event.Block].Mined = true

	for currentNode := 0; currentNode < len(pobSimulation.Nodes); currentNode += 1 {
		if currentNode == event.Node {
			ScheduleBlockMinedEventPoB(currentNode, event.Block)
		} else {
			ScheduleBlockReceivedEventPoB(currentNode, event.Block)
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

func HandleBlockReceivedEventPoB(event *Event) {
	minedBlock := pobSimulation.Nodes[event.Node].CurrentlyMinedBlock

	if pobSimulation.Blocks[minedBlock].PreviousBlock == pobSimulation.Blocks[event.Block].PreviousBlock {
		pobSimulation.Blocks[minedBlock].FinishedAt = event.DispatchAt
		ScheduleBlockMinedEventPoB(event.Node, event.Block)
		return
	}

	if pobSimulation.Blocks[minedBlock].Depth+2 <= pobSimulation.Blocks[event.Block].Depth {
		pobSimulation.Blocks[minedBlock].FinishedAt = event.DispatchAt
		ScheduleBlockMinedEventPoB(event.Node, event.Block)
		return
	}
}
