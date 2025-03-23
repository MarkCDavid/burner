package internal

import "fmt"

func CalculateStatistics(blocks []Block) {

	calculateAverageBlockMiningSpeed(blocks)
	calculateMainChainAverageBlockMiningSpeed(blocks)
	calculateAverageBlockWorkingTime(blocks)

}

func calculateMainChainAverageBlockMiningSpeed(blocks []Block) {
	totalBlocksMined := 0.0
	totalBlockMiningTime := 0.0

	index := getLongestChainTip(blocks)

	for index != -1 {
		if blocks[index].FinishedAt <= 0 {
			index = blocks[index].PreviousBlock
			continue
		}
		totalBlocksMined += 1
		totalBlockMiningTime += (blocks[index].FinishedAt - blocks[index].StartedAt)
		index = blocks[index].PreviousBlock
	}

	fmt.Printf("Average main chain block mining speed: %f\n", totalBlockMiningTime/totalBlocksMined)
}

func calculateAverageBlockMiningSpeed(blocks []Block) {
	totalBlocksMined := 0.0
	totalBlockMiningTime := 0.0

	for _, block := range blocks {
		if !block.Mined {
			continue
		}
		totalBlocksMined += 1
		totalBlockMiningTime += (block.FinishedAt - block.StartedAt)
	}

	fmt.Printf("Average block mining speed: %f\n", totalBlockMiningTime/totalBlocksMined)
}

func calculateAverageBlockWorkingTime(blocks []Block) {
	blockCount := 0.0
	totalWorkTime := 0.0

	for _, block := range blocks {
		if block.FinishedAt <= 0 {
			continue
		}
		blockCount += 1
		totalWorkTime += (block.FinishedAt - block.StartedAt)
	}

	fmt.Printf("Average block work time: %f\n", totalWorkTime/blockCount)
}

func getLongestChainTip(blocks []Block) int {
	max := 0
	index := -1

	for i := range blocks {
		if blocks[i].Depth <= max {
			continue
		}

		index = i
		max = blocks[i].Depth
	}

	return index
}
