package internal

import "github.com/sirupsen/logrus"

func (s *Simulation) GetNextReceivedTime() float64 {
	if s.Configuration.AverageNetworkLatencyInSeconds <= 0 {
		return s.CurrentTime
	}

	return s.CurrentTime + s.Random.Expovariate(1.0/float64(s.Configuration.AverageNetworkLatencyInSeconds))
}

func (s *Simulation) GetTransactionsToInclude(event *Event) int64 {
	availableTransactions := s.GetCurrentTransactionCount() - event.Node.Transactions

	if availableTransactions > s.Configuration.MaximumTransactionsPerBlock {
		return s.Configuration.MaximumTransactionsPerBlock
	}

	if availableTransactions < 0 {
		return 0
	}

	return availableTransactions
}

func ProduceBlock(minedBy *Node, receivedEvent *Event) *Block {
	consensus := minedBy.GetConsensus(receivedEvent)
	if consensus == nil {
		return nil
	}

	minedBy.Simulation.BlockCount += 1

	return &Block{
		Id:        minedBy.Simulation.BlockCount,
		Node:      minedBy,
		Depth:     receivedEvent.Block.Depth + 1,
		Consensus: consensus,
	}
}

func (s *Simulation) ScheduleBlockMinedEvent(
	minedBy *Node,
	receivedEvent *Event,
) {
	block := ProduceBlock(minedBy, receivedEvent)
	if block == nil {
		logrus.Warnf("Miner (%d) can't produce a block in any available consensus layers.", minedBy.Id)
		return
	}

	event := &Event{
		EventType: BlockMinedEvent,

		Block:         block,
		PreviousBlock: receivedEvent.Block,

		ScheduledAt: s.CurrentTime,
	}
	event.SetMiner(minedBy)

	event.DispatchAt = block.Consensus.GetNextMiningTime(event)
	event.Block.Transactions = s.GetTransactionsToInclude(event)

	minedBy.Event = event
	s.Events.Push(event)

}

func (s *Simulation) ScheduleBlockReceivedEvent(receivedBy *Node, minedEvent *Event) {
	e := &Event{
		EventType: BlockReceivedEvent,

		Block:         minedEvent.Block,
		PreviousBlock: minedEvent.PreviousBlock,

		ScheduledAt: s.CurrentTime,
		DispatchAt:  s.GetNextReceivedTime(),
	}
	e.SetReceiver(receivedBy)

	s.Events.Push(e)
}
