package internal

type Event_RandomReceived struct {
	Simulation *Simulation

	Tick int64

	ReceivedBy *Node

	DispatchAt float64

	Index int
}

func (event *Event_RandomReceived) Handle() {
	if event.ReceivedBy.ProofOfBurn == nil {
		return
	}

	event.ReceivedBy.ProofOfBurn.Adjust(event)
	if !event.ReceivedBy.ProofOfBurn.CanMine(event) {
		return
	}

	block := event.ReceivedBy.BuildBlock(event.ReceivedBy.PreviousBlock.Depth+1, event.ReceivedBy.ProofOfBurn)
	event.ReceivedBy.Simulation.ScheduleBlockMinedEvent(event.ReceivedBy, block)
}

func (simulation *Simulation) ScheduleRandomReceivedEvent(receivedBy *Node, tick int64) {
	events := &Event_RandomReceived{
		Simulation: simulation,

		Tick: tick,

		ReceivedBy: receivedBy,

		DispatchAt: simulation.CurrentTime + simulation.GetNextReceivedTime(),
	}
	simulation.Events.Push(events)
}

// === Interface ===

func (e *Event_RandomReceived) GetIndex() int {
	return e.Index
}

func (e *Event_RandomReceived) SetIndex(index int) {
	e.Index = index
}

func (e *Event_RandomReceived) EventTime() float64 {
	return e.DispatchAt
}
