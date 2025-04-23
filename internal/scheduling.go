package internal

import (
	"github.com/sirupsen/logrus"
)

func (s *Simulation) GetNextMiningTime() float64 {
	return s.CurrentTime + s.Random.Expovariate(s.Configuration.InverseAverageBlockFrequency)
}

func (s *Simulation) GetNextReceivedTime() float64 {
	if s.Configuration.AverageNetworkLatencyInSeconds <= 0 {
		return s.CurrentTime
	}

	delta := s.Random.Expovariate(s.Configuration.InverseAverageNetworkLatency)
	if delta > s.Configuration.UpperBoundNetworkLatency {
		delta = s.Configuration.UpperBoundNetworkLatency
	}

	if delta < s.Configuration.LowerBoundNetworkLatency {
		delta = s.Configuration.LowerBoundNetworkLatency
	}

	return s.CurrentTime + delta
}

func (s *Simulation) ScheduleBlockMinedEvent(
	minedBy int,
	previousBlock int,
	depth int,
) {

	s.BlockCount += 1
	logrus.Debugf("SCHEDULING | BlockMinedEvent | Block %d mined by %d (on fork %d)", s.BlockCount, minedBy, s.Nodes[minedBy].Fork)
	minedAt := s.GetNextMiningTime()

	s.Blocks[s.BlockCount] = &Block{
		Node:          minedBy,
		Block:         s.BlockCount,
		PreviousBlock: previousBlock,
		Depth:         depth,
		Fork:          s.Nodes[minedBy].Fork,
		ScheduledAt:   s.CurrentTime,
		DispatchAt:    minedAt,
	}

	if s.Forks[s.Nodes[minedBy].Fork].Len() == 0 {
		logrus.Debugf("EMPTY | Fork %d | Adding a new event.", s.Nodes[minedBy].Fork)
		s.scheduleBlockMinedEvent(minedBy, minedAt, previousBlock, depth)
		return
	}

	nextBlockMinedEvent := s.Forks[s.Nodes[minedBy].Fork].Peek()

	newBlockLowerBound := minedAt
	newBlockUpperBound := minedAt + s.Configuration.UpperBoundNetworkLatency

	oldBlockLowerBound := nextBlockMinedEvent.DispatchAt
	oldBlockUpperBound := nextBlockMinedEvent.DispatchAt + s.Configuration.UpperBoundNetworkLatency

	blockMinedEarlier := newBlockUpperBound < oldBlockLowerBound
	blockMinedLater := newBlockLowerBound > oldBlockUpperBound
	blocksOverlap := !blockMinedEarlier && !blockMinedLater

	if blocksOverlap {
		// The blocks overlap. As such, we will be creating a
		// fork for this. The existing blocks will stay on
		// their own fork, but a new fork will appear.

		s.ForkCount += 1
		s.Forks[s.ForkCount] = CreateEventQueue()
		s.Nodes[minedBy].Fork = s.ForkCount

		logrus.Debugf("OVERLAP | Fork %d | Blocks mined within block receival interval.", s.Nodes[minedBy].Fork)

		s.scheduleBlockMinedEvent(minedBy, minedAt, previousBlock, depth)
		return
	}

	if blockMinedEarlier {
		// The block here appears earlier than the upcoming block.
		// Because of this, we will be removing the existing block
		// and adding the newly scheduled one.

		s.Forks[s.Nodes[minedBy].Fork].Remove(nextBlockMinedEvent)
		s.scheduleBlockMinedEvent(minedBy, minedAt, previousBlock, depth)

		logrus.Debugf("EARLY | Fork %d | Block mined earlier than the upcoming block.", s.Nodes[minedBy].Fork)
		logrus.Debugf("EARLY | Fork %d | Removing %s.", s.Nodes[minedBy].Fork, nextBlockMinedEvent.ToString())

		// TODO: Calculate statistics for the removed block.
		return
	}

	if blockMinedLater {
		// The block here appears later than the upcoming block.
		// Because of this, we will simply calculate the statistics
		// for this block and move on.

		logrus.Debugf("LATE | Fork %d | Block mined later than the upcoming block.", s.Nodes[minedBy].Fork)

		// TODO: Calculate statistics for the removed block.
		// IMPORTANT: Need to calculate receive delay here - across many nodes it matters.
		return
	}

	// if len(s.Forks) > 1 {
	// 	logrus.Debugf("FORKS | Fork %d | Force adding new event.", s.Nodes[minedBy].Fork)
	// 	s.scheduleBlockMinedEvent(minedBy, minedAt, previousBlock, depth)
	// 	return
	// }
	//
	panic("impossible situation: block is not - mined before, mined after, overlaps")
}

func (s *Simulation) scheduleBlockMinedEvent(minedBy int, minedAt float64, previousBlock int, depth int) {

	event := &Event{
		Node: minedBy,

		Block:         s.BlockCount,
		PreviousBlock: previousBlock,
		Depth:         depth,
		Fork:          s.Nodes[minedBy].Fork,

		ScheduledAt: s.CurrentTime,
		DispatchAt:  minedAt,
	}

	s.Nodes[minedBy].CurrentEvent = event
	s.Forks[s.Nodes[minedBy].Fork].Push(event)
}

func (s *Simulation) ScheduleBlockReceivedEvent(receivedBy int, minedEvent *Event) {

	logrus.Debugf("SCHEDULING | BlockReceivedEvent | %s", minedEvent.ToString())
	if len(s.Forks) > 1 {
		event := &Event{
			Node: receivedBy,

			Block:         minedEvent.Block,
			PreviousBlock: minedEvent.PreviousBlock,
			Depth:         minedEvent.Depth,
			Fork:          minedEvent.Fork,

			ScheduledAt: s.CurrentTime,
			DispatchAt:  s.GetNextReceivedTime(),
		}

		logrus.Debugf("SCHEDULING | BlockReceivedEvent | Block %d received by %d | %s", minedEvent.Block, receivedBy, event.ToString())

		s.ForkDependence[minedEvent.Fork] += 1
		s.Network.Push(event)
		return
	}

	oldSimulationTime := s.CurrentTime
	s.CurrentTime = s.GetNextReceivedTime()

	logrus.Debugf("SCHEDULING | BlockReceivedEvent | Block %d received by %d | Pretending that current time is %f", minedEvent.Block, receivedBy, s.CurrentTime)
	s.Nodes[receivedBy].Fork = minedEvent.Fork
	s.ScheduleBlockMinedEvent(receivedBy, minedEvent.Block, minedEvent.Depth+1)
	s.CurrentTime = oldSimulationTime
}
