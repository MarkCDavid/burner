package internal

type NodeStatistics struct {
	BlocksMined               [2]int64
	TransactionsProcessed     [2]int64
	BlockMiningTime           [2]float64
	BlockSuccessfulMiningTime [2]float64

	PowerUsed [2]float64
}

type Statistics struct {
	BlocksMined           [2]int64
	TransactionsProcessed [2]int64
	BlockMiningTime       [2]float64
	BlockIntervalTime     [2]float64
	PerNode               []NodeStatistics

	ForkResolutions int64
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

func (s *Statistics) OnBlockMined(simulation *Simulation, event *Event_BlockMined) {
	simulation.Database.SaveBlock(event)

	node := event.MinedBy.Id
	consensusType := event.Block.Consensus.GetType()

	s.BlocksMined[consensusType] += 1
	s.TransactionsProcessed[consensusType] += event.Block.Transactions
	s.BlockMiningTime[consensusType] += event.Block.FinishedAt - event.Block.StartedAt
	s.BlockIntervalTime[consensusType] += event.Block.FinishedAt - event.PreviousBlock.FinishedAt

	s.PerNode[node].BlocksMined[consensusType] += 1
	s.PerNode[node].PowerUsed[consensusType] += event.PowerUsed()
	s.PerNode[node].TransactionsProcessed[consensusType] += event.Block.Transactions

	s.PerNode[node].BlockMiningTime[consensusType] += event.Duration()
	s.PerNode[node].BlockSuccessfulMiningTime[consensusType] += event.Duration()
}

func (s *Statistics) OnBlockAbandoned(simulation *Simulation, event *Event_BlockMined) {
	simulation.Database.SaveBlock(event)

	node := event.MinedBy.Id
	consensusType := event.Block.Consensus.GetType()

	s.PerNode[node].PowerUsed[consensusType] += event.PowerUsed()
	s.PerNode[node].BlockMiningTime[consensusType] += event.Duration()
}

func (s *Statistics) GetTotalBlocks() int64 {
	return s.BlocksMined[ProofOfWork] + s.BlocksMined[ProofOfBurn]
}

func (s *Statistics) GetTotalTransactions() int64 {
	return s.TransactionsProcessed[ProofOfWork] + s.TransactionsProcessed[ProofOfBurn]
}

func (s *Statistics) GetTotalMiningTime() float64 {
	return s.BlockMiningTime[ProofOfWork] + s.BlockMiningTime[ProofOfBurn]
}

func (s *Statistics) GetTotalPowerUsed() float64 {
	var total float64 = 0
	for i := 0; i < len(s.PerNode); i++ {
		total += s.PerNode[i].PowerUsed[ProofOfWork]
		total += s.PerNode[i].PowerUsed[ProofOfBurn]
	}
	return total
}
