package internal

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

func (s *Simulation) ScheduleBlockMinedEvent(
	minedBy int64,
	receivedEvent *Event,
) {
	s.BlockCount += 1

	event := &Event{
		Node:      minedBy,
		EventType: BlockMinedEvent,

		Block: &Block{
			Id:    s.BlockCount,
			Node:  minedBy,
			Depth: receivedEvent.Block.Depth + 1,
			Type:  ProofOfWork,
		},
		PreviousBlock: receivedEvent.Block,

		ScheduledAt: s.CurrentTime,
	}

	event.DispatchAt = s.GetNextMiningTime(event)

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
