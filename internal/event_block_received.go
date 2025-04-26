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

func (e *Event_BlockReceived) SetReceiver(n *Node) {
	e.ReceivedBy = n
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
		event.Reorganize()
		return
	}
}

// TODO: Work on me please
func (event *Event_BlockReceived) Reorganize() {
	if event.ReceivedBy != event.Block.Node {
		event.ReceivedBy.SynchronizeConsensus(event.Block.Node)
		event.ReceivedBy.Transactions = event.Block.Node.Transactions

		if event.ReceivedBy.Event != nil {
			event.ReceivedBy.Simulation.Statistics.OnBlockAbandoned(event.Simulation, event.ReceivedBy.Event)
			event.Simulation.Events.Remove(event.ReceivedBy.Event)
		}
	}

	event.Simulation.ScheduleBlockMinedEvent(event.ReceivedBy, event)
}

func (s *Simulation) GetNextReceivedTime() float64 {
	if s.Configuration.AverageNetworkLatencyInSeconds <= 0 {
		return s.CurrentTime
	}

	return s.CurrentTime + s.Random.Expovariate(1.0/float64(s.Configuration.AverageNetworkLatencyInSeconds))
}

func (s *Simulation) ScheduleBlockReceivedEvent(receivedBy *Node, event *Event_BlockMined) {
	receivedTime := event.EventTime()
	if receivedBy != event.MinedBy {
		receivedTime = s.GetNextReceivedTime()
	}

	e := &Event_BlockReceived{
		Simulation:    s,
		Block:         event.Block,
		PreviousBlock: event.PreviousBlock,

		ScheduledAt: s.CurrentTime,
		DispatchAt:  receivedTime,
	}
	e.SetReceiver(receivedBy)
	s.Events.Push(e)
}

func (e *Event_BlockReceived) Duration() float64 {
	return e.Simulation.CurrentTime - e.ScheduledAt
}

func (event *Event_BlockReceived) GetIndex() int {
	return event.Index
}

func (event *Event_BlockReceived) SetIndex(index int) {
	event.Index = index
}

func (event *Event_BlockReceived) EventTime() float64 {
	return event.DispatchAt
}
