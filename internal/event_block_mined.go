package internal

// import "github.com/sirupsen/logrus"

type Event_BlockMined struct {
	Simulation *Simulation

	MinedBy *Node

	Block         *Block
	PreviousBlock *Block

	ScheduledAt float64
	DispatchAt  float64

	Index int
}

func (event *Event_BlockMined) Handle() {
	event.Block.FinishedAt = event.Simulation.CurrentTime
	event.MinedBy.PreviousBlock = event.Block
	// logrus.Infof("%d mined %s block (%d <- %d)", event.MinedBy.Id, event.Block.Consensus.GetType().ToString(), event.PreviousBlock.Id, event.Block.Id)
	event.Simulation.Statistics.OnBlockMined(event.MinedBy.Simulation, event)

	if event.MinedBy.ProofOfWork != nil {
		event.MinedBy.ProofOfWork.Adjust(event)
	}

	if event.MinedBy.ProofOfBurn != nil {
		event.MinedBy.ProofOfBurn.Adjust(event)
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

func (minedBy *Node) BuildBlock(depth int64, consensus Consensus) *Block {
	minedBy.Simulation.BlockCount += 1

	return &Block{
		Id:        minedBy.Simulation.BlockCount,
		Node:      minedBy,
		Depth:     depth,
		Consensus: consensus,
	}
}

func (minedBy *Node) ProduceBlock(event *Event_BlockReceived) *Block {
	consensus := minedBy.GetConsensus(event)
	if consensus == nil {
		return nil
	}

	return minedBy.BuildBlock(event.Block.Depth+1, consensus)
}

func (simulation *Simulation) ScheduleBlockMinedEvent(
	minedBy *Node,
	block *Block,
) {
	if block == nil {
		return
	}

	event := &Event_BlockMined{
		Simulation: simulation,

		MinedBy: minedBy,

		Block:         block,
		PreviousBlock: minedBy.PreviousBlock,

		ScheduledAt: simulation.CurrentTime,
	}
	minedBy.Event = event
	event.DispatchAt = block.Consensus.GetNextMiningTime(event)
	event.Block.Transactions = simulation.GetTransactionsToInclude(event)

	simulation.Events.Push(event)
}

func (e *Event_BlockMined) PowerUsed() float64 {
	return e.Duration() * e.Block.Consensus.GetPower()
}

func (e *Event_BlockMined) Duration() float64 {
	return e.Simulation.CurrentTime - e.PreviousBlock.FinishedAt
}

func (e *Event_BlockMined) GetIndex() int {
	return e.Index
}

func (e *Event_BlockMined) SetIndex(index int) {
	e.Index = index
}

func (e *Event_BlockMined) EventTime() float64 {
	return e.DispatchAt
}
