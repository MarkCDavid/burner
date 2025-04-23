package internal

func (s *Simulation) GetNextMiningTime() float64 {
	return s.CurrentTime + s.Random.Expovariate(s.Configuration.InverseAverageBlockFrequency)
}

func (s *Simulation) GetNextReceivedTime() float64 {
	if s.Configuration.AverageNetworkLatencyInSeconds <= 0 {
		return s.CurrentTime
	}

	return s.CurrentTime + s.Random.Expovariate(s.Configuration.InverseAverageNetworkLatency)
}

func (s *Simulation) ScheduleBlockMinedEvent(
	minedBy int,
	previousBlock int,
	depth int,
) {
	s.BlockCount += 1
	event := &Event{
		Node:      minedBy,
		EventType: BlockMinedEvent,

		Block:         s.BlockCount,
		PreviousBlock: previousBlock,
		Depth:         depth,

		ScheduledAt: s.CurrentTime,
		DispatchAt:  s.GetNextMiningTime(),
	}

	s.Nodes[minedBy].CurrentEvent = event
	s.Events.Push(event)
}

func (s *Simulation) ScheduleBlockReceivedEvent(receivedBy int, minedEvent *Event) {
	event := &Event{
		Node:      receivedBy,
		EventType: BlockReceivedEvent,

		Block:         minedEvent.Block,
		PreviousBlock: minedEvent.PreviousBlock,
		Depth:         minedEvent.Depth,

		ScheduledAt: s.CurrentTime,
		DispatchAt:  s.GetNextReceivedTime(),
	}

	s.Events.Push(event)
}
