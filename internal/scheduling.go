package internal

import "github.com/sirupsen/logrus"

func (s *Simulation) GetNextMiningTime(event *Event) float64 {
	lambda := s.Nodes[event.Node].Difficulty[event.Block.Type].GetLambda(s.Nodes[event.Node].Power)
	return s.CurrentTime + s.Random.Expovariate(lambda)
}

func (s *Simulation) GetNextReceivedTime() float64 {
	if s.Configuration.AverageNetworkLatencyInSeconds <= 0 {
		return s.CurrentTime
	}

	return s.CurrentTime + s.Random.Expovariate(1.0/float64(s.Configuration.AverageNetworkLatencyInSeconds))
}

func (s *Simulation) GetTransactionsToInclude(event *Event) int64 {
	availableTransactions := s.GetCurrentTransactionCount() - s.Nodes[event.Node].Transactions

	if availableTransactions > s.Configuration.MaximumTransactionsPerBlock {
		return s.Configuration.MaximumTransactionsPerBlock
	}

	if availableTransactions < 0 {
		return 0
	}

	return availableTransactions
}

func (s *Simulation) ScheduleBlockMinedEvent(
	minedBy int64,
	receivedEvent *Event,
) {

	var blockType BlockType = DifficultyVariants - 1

	for {
		if s.Nodes[minedBy].Difficulty[blockType].CanMine(s, receivedEvent.Block, s.Nodes[minedBy].Power) {
			break
		}

		blockType--
		if blockType < 0 {
			logrus.Warnf("%d - miner could not mine a single block type", minedBy)
			return
		}
	}

	s.BlockCount += 1

	event := &Event{
		Node:      minedBy,
		EventType: BlockMinedEvent,

		Block: &Block{
			Id:    s.BlockCount,
			Node:  minedBy,
			Depth: receivedEvent.Block.Depth + 1,
			Type:  blockType,
		},
		PreviousBlock: receivedEvent.Block,

		ScheduledAt: s.CurrentTime,
	}

	event.DispatchAt = s.GetNextMiningTime(event)
	event.Block.Transactions = s.GetTransactionsToInclude(event)

	s.Nodes[minedBy].CurrentEvent = event
	s.Events.Push(event)
}

func (s *Simulation) ScheduleBlockReceivedEvent(receivedBy int64, minedEvent *Event) {
	event := &Event{
		Node:      receivedBy,
		EventType: BlockReceivedEvent,

		Block:         minedEvent.Block,
		PreviousBlock: minedEvent.PreviousBlock,

		ScheduledAt: s.CurrentTime,
		DispatchAt:  s.GetNextReceivedTime(),
	}

	s.Events.Push(event)
}
