package internal

type Statistics struct {
	MinedBlockCount [2]int64
	BlockMiningTime [2]float64
}

func (s *Statistics) GetAverageBlockMiningTime() float64 {
	var totalBlocksMined int64 = 0
	var totalMiningTime float64 = 0
	for i := 0; i < len(s.BlockMiningTime); i++ {
		totalBlocksMined += s.MinedBlockCount[i]
		totalMiningTime += s.BlockMiningTime[i]
	}
	return totalMiningTime / float64(totalBlocksMined)
}

func (s *Statistics) OnBlockMined(event *Event) {
	s.MinedBlockCount[event.Block.Type] += 1
	s.BlockMiningTime[event.Block.Type] += event.Duration()
}
