package internal

import (
	"fmt"
)

func GetNextMiningTime(nodeIndex int, previousBlock int) float64 {
	difficulty := simulation.Blocks[previousBlock].ProofOfBurnDifficulty
	distributedFrequency := float64(simulation.Configuration.AverageBlockFrequencyInSeconds) * float64(difficulty)
	return simulation.CurrentTime + simulation.Random.Expovariate((simulation.Nodes[nodeIndex].NodePower)/distributedFrequency)
}

func GetNextMiningTimePoB(nodeIndex int, previousBlock int) float64 {
	// TODO: Need to model when the "virtual pc" cannot mine any longer.
	difficulty := pobSimulation.Blocks[previousBlock].ProofOfBurnDifficulty
	distributedFrequency := float64(600) * float64(difficulty)
	return pobSimulation.CurrentTime + pobSimulation.Random.Expovariate(100/distributedFrequency)
}
func GetNextReceivedTimePoB() float64 {
	offset := pobSimulation.Random.Expovariate(1.0 / float64(6))
	return pobSimulation.CurrentTime + offset
}

func GetNextReceivedTime() float64 {
	if simulation.Configuration.AverageNetworkLatencyInSeconds <= 0 {
		return simulation.CurrentTime
	}

	offset := simulation.Random.Expovariate(1.0 / float64(simulation.Configuration.AverageNetworkLatencyInSeconds))

	if offset >= 2*simulation.Configuration.SimulationTime {
		offset = 2 * simulation.Configuration.SimulationTime
	}

	if offset <= 0.5*simulation.Configuration.SimulationTime {
		offset = 0.5 * simulation.Configuration.SimulationTime
	}

	return simulation.CurrentTime + offset
}

func GetNextBlockType(minedBy int) BlockType {
	return ProofOfWork
	// if simulation.Configuration.ProofOfBurnEveryNthBlock <= 0 {
	// 	return ProofOfBurn
	// }
	// proofOfBurnProbability := 1.0 / float64(simulation.Configuration.NodeCount*simulation.Configuration.ProofOfBurnEveryNthBlock)
	// if proofOfBurnProbability > simulation.Random.float() {
	// 	return ProofOfBurn
	// }
	//
	// return ProofOfWork
}
func ScheduleBlockMinedEventPoB(
	minedBy int,
	previousBlock int,
) {
	difficulty := pobSimulation.Blocks[previousBlock].ProofOfBurnDifficulty
	minedAt := GetNextMiningTimePoB(minedBy, previousBlock)

	miningTime := (minedAt - pobSimulation.CurrentTime)

	difficultyMultiplier := 600.0 / miningTime
	if difficultyMultiplier > 4 {
		difficultyMultiplier = 4
	}
	if difficultyMultiplier < 0.25 {
		difficultyMultiplier = 0.25
	}

	block := Block{
		Node:                  minedBy,
		BlockType:             ProofOfBurn,
		PreviousBlock:         previousBlock,
		StartedAt:             pobSimulation.CurrentTime,
		FinishedAt:            minedAt,
		ProofOfBurnDifficulty: difficulty * difficultyMultiplier,
		Depth:                 pobSimulation.Blocks[previousBlock].Depth + 1,
	}

	pobSimulation.Blocks = append(pobSimulation.Blocks, block)
	currentBlock := len(pobSimulation.Blocks) - 1
	pobSimulation.Nodes[minedBy].CurrentlyMinedBlock = currentBlock

	pobSimulation.Queue.Push(&Event{
		Type:       BlockMinedEvent,
		Node:       minedBy,
		Block:      currentBlock,
		DispatchAt: minedAt,
	})
}

func ScheduleBlockMinedEvent(
	minedBy int,
	previousBlock int,
) {
	blockType := GetNextBlockType(minedBy)
	minedAt := GetNextReceivedTime()
	difficulty := simulation.Blocks[previousBlock].ProofOfBurnDifficulty
	finishedAt := minedAt

	// nextEvent := simulation.Queue.Peek()
	// if nextEvent.Type == BlockMinedEvent && nextEvent.DispatchAt+2*simulation.Configuration.SimulationTime < minedAt {
	//
	// }
	//
	switch blockType {
	case ProofOfBurn:
		break
	case ProofOfWork:
		minedAt = GetNextMiningTime(minedBy, previousBlock)

		if simulation.Blocks[previousBlock].Depth > 0 && simulation.Blocks[previousBlock].Depth%2016 == 0 {
			//Difficulty adjustment TIME!!!!
			currentIndex := previousBlock
			previousIndex := simulation.Blocks[currentIndex].PreviousBlock
			totalTime := 0.0
			for i := 0; i < 2016; i++ {

				totalTime += simulation.Blocks[currentIndex].FinishedAt - simulation.Blocks[previousIndex].FinishedAt
				currentIndex = previousIndex
				previousIndex = simulation.Blocks[currentIndex].PreviousBlock
			}
			fmt.Printf("Epoch time %f\n", totalTime)
			fmt.Printf("Expected ~2 weeks %f\n", 2.0*7.0*24.0*60.0*60.0)
			fmt.Printf("Diff %f\n", totalTime-2.0*7.0*24.0*60.0*60.0)

			difficultyMultiplier := (2.0 * 7.0 * 24.0 * 60.0 * 60.0) / float64(totalTime)
			if difficultyMultiplier < 0.25 {
				difficultyMultiplier = 0.25
			}
			if difficultyMultiplier > 4 {
				difficultyMultiplier = 4
			}

			difficulty = difficulty * difficultyMultiplier
		}

		finishedAt = 0
		break
	default:
		panic("Unknown type")
	}

	block := Block{
		Node:                  minedBy,
		BlockType:             blockType,
		PreviousBlock:         previousBlock,
		StartedAt:             simulation.CurrentTime,
		FinishedAt:            finishedAt,
		ProofOfBurnDifficulty: difficulty,
		Depth:                 simulation.Blocks[previousBlock].Depth + 1,
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
func ScheduleBlockReceivedEventPoB(receivedBy int, minedBlock int) {
	pobSimulation.Queue.Push(&Event{
		Type:       BlockReceivedEvent,
		Node:       receivedBy,
		Block:      minedBlock,
		DispatchAt: GetNextReceivedTimePoB(),
	})
}
