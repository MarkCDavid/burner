package internal

func GetNextMiningTime() float64 {
	distributedFrequency := simulation.Configuration.AverageBlockFrequencyInSeconds * simulation.Configuration.NodeCount
	return simulation.CurrentTime + simulation.Random.Expovariate(1.0/float64(distributedFrequency))
}

func GetNextReceivedTime() float64 {
	if simulation.Configuration.AverageNetworkLatencyInSeconds <= 0 {
		return simulation.CurrentTime
	}

	return simulation.CurrentTime + simulation.Random.Expovariate(1.0/float64(simulation.Configuration.AverageNetworkLatencyInSeconds))
}

func GetNextBlockType() BlockType {
	if simulation.Configuration.ProofOfBurnEveryNthBlock <= 0 {
		return ProofOfBurn
	}
	proofOfBurnProbability := 1.0 / float64(simulation.Configuration.NodeCount*simulation.Configuration.ProofOfBurnEveryNthBlock)
	if proofOfBurnProbability > simulation.Random.float() {
		return ProofOfBurn
	}

	return ProofOfWork
}

func ScheduleBlockMinedEvent(
	minedBy int,
	previousBlock int,
) {
	blockType := GetNextBlockType()
	minedAt := GetNextReceivedTime()
	finishedAt := minedAt
	switch blockType {
	case ProofOfBurn:
		break
	case ProofOfWork:
		minedAt = GetNextMiningTime()
		finishedAt = 0
		break
	default:
		panic("Unknown type")
	}

	block := Block{
		Node:          minedBy,
		BlockType:     blockType,
		PreviousBlock: previousBlock,
		StartedAt:     simulation.CurrentTime,
		FinishedAt:    finishedAt,
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

func ScheduleBlockReceivedEvent(receivedBy int, minedBlock int) {
	simulation.Queue.Push(&Event{
		Type:       BlockReceivedEvent,
		Node:       receivedBy,
		Block:      minedBlock,
		DispatchAt: GetNextReceivedTime(),
	})
}
