package internal

type NodeStatistics struct {
	BlocksMined           int64
	TransactionsProcessed int64
	PowerUsed             float64

	TimeMining             float64
	TimeOnSuccessfulMining float64
}

type Statistics struct {
	BlocksMined           [2]int64
	TransactionsProcessed [2]int64
	BlockMiningTime       [2]float64
	PerNode               []NodeStatistics
}

func (s *Statistics) GetAverageBlockMiningTime() float64 {
	var totalBlocksMined int64 = 0
	var totalMiningTime float64 = 0
	for i := 0; i < len(s.BlockMiningTime); i++ {
		totalBlocksMined += s.BlocksMined[i]
		totalMiningTime += s.BlockMiningTime[i]
	}
	return totalMiningTime / float64(s.GetTotalBlocks())
}

func (s *Statistics) GetAverageTransactionsPerBlock() float64 {
	return float64(s.GetTotalTransactions()) / float64(s.GetTotalBlocks())
}

func (s *Statistics) OnBlockMined(simulation *Simulation, event *Event) {
	node := event.Node.Id
	consensusType := event.Block.Consensus.GetType()

	s.BlocksMined[consensusType] += 1
	s.TransactionsProcessed[consensusType] += event.Block.Transactions
	s.BlockMiningTime[consensusType] += event.Duration()

	s.PerNode[node].BlocksMined += 1
	s.PerNode[node].PowerUsed += event.Duration() * simulation.Nodes[node].Capability
	s.PerNode[node].TransactionsProcessed += event.Block.Transactions

	s.PerNode[node].TimeMining += event.Duration()
	s.PerNode[node].TimeOnSuccessfulMining += event.Duration()
}

func (s *Statistics) OnBlockAbandoned(simulation *Simulation, event *Event) {
	node := event.Node.Id
	s.PerNode[node].PowerUsed += event.Duration() * simulation.Nodes[node].Capability
	s.PerNode[node].TimeMining += event.Duration()
}

func (s *Statistics) GetTotalBlocks() int64 {
	var total int64 = 0
	for i := 0; i < DifficultyVariants; i++ {
		total += s.BlocksMined[i]
	}
	return total
}

func (s *Statistics) GetTotalTransactions() int64 {
	var total int64 = 0
	for i := 0; i < DifficultyVariants; i++ {
		total += s.TransactionsProcessed[i]
	}
	return total
}

func (s *Statistics) GetTotalMiningTime() float64 {
	var total float64 = 0
	for i := 0; i < DifficultyVariants; i++ {
		total += s.BlockMiningTime[i]
	}
	return total
}

func (s *Statistics) GetTotalPowerUsed() float64 {
	var total float64 = 0
	for i := 0; i < len(s.PerNode); i++ {
		total += s.PerNode[i].PowerUsed
	}
	return total
}
