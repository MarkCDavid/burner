package internal

func (s *Simulation) HandleBlockMinedEvent(event *Event) {
	s.Statistics.OnBlockMined(s, event)

	event.Block.Consensus.Adjust(event)
	event.Node.Transactions += event.Block.Transactions

	s.ScheduleBlockMinedEvent(event.Node, event)

	for _, node := range s.Nodes {
		if node != event.Node {
			s.ScheduleBlockReceivedEvent(node, event)
		}
	}
}

func (s *Simulation) HandleBlockReceivedEvent(event *Event) {
	if event.Node.Event.PreviousBlock.Id == event.PreviousBlock.Id {
		s.Reorganize(event)
		return
	}

	reorganizationThreshold := event.Node.Event.Block.Depth + s.Configuration.ChainReogranizationThreshold
	deepEnough := reorganizationThreshold <= event.Block.Depth

	if deepEnough {
		s.Reorganize(event)
		return
	}
}

func (s *Simulation) Reorganize(event *Event) {
	event.Node.SynchronizeConsensus(event.Block.Node)
	event.Node.Transactions = event.Block.Node.Transactions

	s.Statistics.OnBlockAbandoned(s, event.Node.Event)
	s.Events.Remove(event.Node.Event)

	s.ScheduleBlockMinedEvent(event.Node, event)
}
