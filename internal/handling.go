package internal

func (s *Simulation) HandleBlockMinedEvent(event *Event) {
	s.Statistics.OnBlockMined(s, event)
	s.Nodes[event.Node].Difficulty[event.Block.Type].Adjust(event)
	s.Nodes[event.Node].Transactions += event.Block.Transactions

	s.ScheduleBlockMinedEvent(event.Node, event)
	for currentNode := int64(0); currentNode < int64(len(s.Nodes)); currentNode += 1 {
		if currentNode != event.Node {
			s.ScheduleBlockReceivedEvent(currentNode, event)
		}
	}
}

func (s *Simulation) HandleBlockReceivedEvent(event *Event) {
	miningEvent := s.Nodes[event.Node].CurrentEvent

	if miningEvent.PreviousBlock.Id == event.PreviousBlock.Id {
		s.Nodes[event.Node].Difficulty[event.Block.Type].Update(s.Nodes[event.Block.Node].Difficulty[event.Block.Type])
		s.Nodes[event.Node].Transactions = s.Nodes[event.Block.Node].Transactions
		s.Statistics.OnBlockAbandoned(s, miningEvent)
		s.Events.Remove(miningEvent)
		s.ScheduleBlockMinedEvent(event.Node, event)
		return
	}

	chainReorganizationThreshold := miningEvent.Block.Depth + s.Configuration.ChainReogranizationThreshold
	deepEnoughForChainReorganization := chainReorganizationThreshold <= event.Block.Depth

	if deepEnoughForChainReorganization {
		s.Nodes[event.Node].Difficulty[event.Block.Type].Update(s.Nodes[event.Block.Node].Difficulty[event.Block.Type])
		s.Nodes[event.Node].Transactions = s.Nodes[event.Block.Node].Transactions
		s.Statistics.OnBlockAbandoned(s, miningEvent)
		s.Events.Remove(miningEvent)
		s.ScheduleBlockMinedEvent(event.Node, event)
		return
	}
}
