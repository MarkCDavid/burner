package internal

type Event_BlockMined struct {
	MinedBy    *Node
	Simulation *Simulation

	Block         *Block
	PreviousBlock *Block

	ScheduledAt float64
	DispatchAt  float64

	Index int
}

func (event *Event_BlockMined) Handle() {
	event.Simulation.Statistics.OnBlockMined(event.MinedBy.Simulation, event)

	for _, consensus := range event.MinedBy.Consensus {
		consensus.Adjust(event)
	}
	event.MinedBy.Transactions += event.Block.Transactions

	for _, node := range event.MinedBy.Simulation.Nodes {
		event.Simulation.ScheduleBlockReceivedEvent(node, event)
	}
}

func (s *Simulation) GetTransactionsToInclude(event *Event_BlockMined) int64 {
	availableTransactions := s.GetCurrentTransactionCount() - event.MinedBy.Transactions

	if availableTransactions > s.Configuration.MaximumTransactionsPerBlock {
		return s.Configuration.MaximumTransactionsPerBlock
	}

	if availableTransactions < 0 {
		return 0
	}

	return availableTransactions
}

func ProduceBlock(minedBy *Node, event *Event_BlockReceived) *Block {
	consensus := minedBy.GetConsensus(event)
	if consensus == nil {
		return nil
	}

	minedBy.Simulation.BlockCount += 1

	return &Block{
		Id:        minedBy.Simulation.BlockCount,
		Node:      minedBy,
		Depth:     event.Block.Depth + 1,
		Consensus: consensus,
	}
}

func (s *Simulation) ScheduleBlockMinedEvent(
	minedBy *Node,
	receivedEvent *Event_BlockReceived,
) {
	block := ProduceBlock(minedBy, receivedEvent)
	if block == nil {
		return
	}

	event := &Event_BlockMined{
		Simulation:    s,
		Block:         block,
		PreviousBlock: receivedEvent.Block,

		ScheduledAt: s.CurrentTime,
	}
	event.SetMiner(minedBy)

	event.DispatchAt = block.Consensus.GetNextMiningTime(event)
	event.Block.Transactions = s.GetTransactionsToInclude(event)

	s.Events.Push(event)
}

func (e *Event_BlockMined) PowerUsed() float64 {
	return e.Duration() * e.MinedBy.Power[e.Block.Consensus.GetType()]
}
func (e *Event_BlockMined) SetMiner(n *Node) {
	e.MinedBy = n
	n.Event = e
}

func (e *Event_BlockMined) Duration() float64 {
	return e.Simulation.CurrentTime - e.ScheduledAt
}

func (event *Event_BlockMined) GetIndex() int {
	return event.Index
}

func (event *Event_BlockMined) SetIndex(index int) {
	event.Index = index
}

func (event *Event_BlockMined) EventTime() float64 {
	return event.DispatchAt
}
