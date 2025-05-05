package internal

type Event_BlockReceived struct {
	ReceivedBy *Node
	Simulation *Simulation

	Block         *Block
	PreviousBlock *Block

	ScheduledAt float64
	DispatchAt  float64

	Index int
}

func (event *Event_BlockReceived) Handle() {
	if event.ReceivedBy.Event == nil {
		event.Reorganize()
		return
	}

	if event.ReceivedBy.Event.PreviousBlock.Id == event.PreviousBlock.Id {
		event.Reorganize()
		return
	}

	reorganizationThreshold := event.ReceivedBy.Event.Block.Depth + event.Simulation.Configuration.ChainReogranizationThreshold
	deepEnough := reorganizationThreshold <= event.Block.Depth

	if deepEnough {
		event.Simulation.Statistics.ForkResolutions++
		event.Reorganize()
		return
	}
}

func (event *Event_BlockReceived) Reorganize() {
	if event.Block.Node != nil && event.ReceivedBy != event.Block.Node {
		event.ReceivedBy.SynchronizeConsensus(event.Block.Node)
		event.ReceivedBy.Transactions = event.Block.Node.Transactions

		if event.ReceivedBy.Event != nil {
			event.ReceivedBy.Event.Block.Abandoned = true
			event.ReceivedBy.Event.Block.FinishedAt = event.ReceivedBy.Simulation.CurrentTime
			event.ReceivedBy.Simulation.Statistics.OnBlockAbandoned(event.Simulation, event.ReceivedBy.Event)
			event.Simulation.Events.Remove(event.ReceivedBy.Event)
		}
	}

	event.ReceivedBy.PreviousBlock = event.Block

	if event.ReceivedBy.ProofOfWork != nil {
		event.ReceivedBy.ProofOfWork.Adjust(event)
	}

	if event.ReceivedBy.ProofOfBurn != nil {
		event.ReceivedBy.ProofOfBurn.Adjust(event)
	}

	block := event.ReceivedBy.ProduceBlock(event)
	event.Simulation.ScheduleBlockMinedEvent(event.ReceivedBy, block)
}

func (s *Simulation) GetNextReceivedTime() float64 {
	if s.Configuration.AverageNetworkLatencyInSeconds <= 0 {
		return s.CurrentTime
	}

	return s.CurrentTime + s.Random.Expovariate(1.0/float64(s.Configuration.AverageNetworkLatencyInSeconds))
}

func (simulation *Simulation) ScheduleBlockReceivedEvent(receivedBy *Node, event *Event_BlockMined) {
	receivedTime := event.EventTime()
	if receivedBy != event.MinedBy {
		receivedTime = simulation.GetNextReceivedTime()
	}

	e := &Event_BlockReceived{
		Simulation: simulation,

		ReceivedBy: receivedBy,

		Block:         event.Block,
		PreviousBlock: event.PreviousBlock,

		ScheduledAt: simulation.CurrentTime,
		DispatchAt:  receivedTime,
	}
	simulation.Events.Push(e)
}

func (e *Event_BlockReceived) Duration() float64 {
	return e.DispatchAt - e.ScheduledAt
}

func (e *Event_BlockReceived) GetIndex() int {
	return e.Index
}

func (e *Event_BlockReceived) SetIndex(index int) {
	e.Index = index
}

func (e *Event_BlockReceived) EventTime() float64 {
	return e.DispatchAt
}
