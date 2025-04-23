package internal

func (s *Simulation) HandleBlockMinedEvent(event *Event) {
	s.ScheduleBlockMinedEvent(event.Node, event.Block, event.Depth+1)
	for currentNode := 0; currentNode < len(s.Nodes); currentNode += 1 {
		if currentNode != event.Node {
			s.ScheduleBlockReceivedEvent(currentNode, event)
		}
	}
}

func (s *Simulation) HandleBlockReceivedEvent(event *Event) {
	miningEvent := s.Nodes[event.Node].CurrentEvent

	if miningEvent.PreviousBlock == event.PreviousBlock {
		s.Events.Remove(miningEvent)
		s.ScheduleBlockMinedEvent(event.Node, event.Block, event.Depth+1)
		return
	}

	chainReorganizationThreshold := miningEvent.Depth + s.Configuration.ChainReogranizationThreshold
	deepEnoughForChainReorganization := chainReorganizationThreshold <= event.Depth

	if deepEnoughForChainReorganization {
		s.Events.Remove(miningEvent)
		s.ScheduleBlockMinedEvent(event.Node, event.Block, event.Depth+1)
		return
	}
}
