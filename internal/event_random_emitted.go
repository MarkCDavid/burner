package internal

type Event_RandomEmitted struct {
	Simulation *Simulation

	Tick int64

	Delay      float64
	DispatchAt float64

	Index int
}

func (event *Event_RandomEmitted) Handle() {
	event.Simulation.ScheduleEmitRandomEvent(event.Tick+1, event.Delay)

	for _, node := range event.Simulation.Nodes {
		event.Simulation.ScheduleRandomReceivedEvent(node, event.Tick)
	}
}

func (simulation *Simulation) ScheduleEmitRandomEvent(tick int64, delay float64) {
	events := &Event_RandomEmitted{
		Simulation: simulation,

		Tick: tick,

		Delay:      delay,
		DispatchAt: simulation.CurrentTime + delay,
	}
	simulation.Events.Push(events)
}

// === Interface ===

func (e *Event_RandomEmitted) GetIndex() int {
	return e.Index
}

func (e *Event_RandomEmitted) SetIndex(index int) {
	e.Index = index
}

func (e *Event_RandomEmitted) EventTime() float64 {
	return e.DispatchAt
}
