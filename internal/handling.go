package internal

import (
	"github.com/sirupsen/logrus"
)

func (s *Simulation) HandleBlockMinedEvent(event *Event) {
	s.Blocks[event.Block].Mined = true
	s.Blocks[event.Block].Fork = s.Nodes[event.Node].Fork
	s.ScheduleBlockMinedEvent(event.Node, event.Block, event.Depth+1)
	for currentNode := 0; currentNode < len(s.Nodes); currentNode += 1 {
		if currentNode != event.Node {
			s.ScheduleBlockReceivedEvent(currentNode, event)
		}
	}
}

func (s *Simulation) HandleBlockReceivedEvent(event *Event) {
	miningEvent := s.Nodes[event.Node].CurrentEvent

	if miningEvent == nil {
		s.Nodes[event.Node].Fork = event.Fork
		s.ScheduleBlockMinedEvent(event.Node, event.Block, event.Depth+1)
		return
	}

	s.ForkDependence[event.Fork] -= 1

	logrus.Debugf("BLOCK RECEIVED | Block Received Event | %s", event.ToString())
	logrus.Debugf("BLOCK RECEIVED | Current Mining Event | %s", miningEvent.ToString())

	if miningEvent.Fork == event.Fork && miningEvent.PreviousBlock == event.PreviousBlock {
		logrus.Debug("BLOCK RECEIVED | rescheduling")

		s.Forks[s.Nodes[event.Node].Fork].Remove(miningEvent)
		s.ScheduleBlockMinedEvent(event.Node, event.Block, event.Depth+1)
		return
	}

	chainReorganizationThreshold := miningEvent.Depth + s.Configuration.ChainReogranizationThreshold
	deepEnoughForChainReorganization := chainReorganizationThreshold <= event.Depth

	if deepEnoughForChainReorganization {
		logrus.Debug("BLOCK RECEIVED | reoganizing")

		s.Forks[s.Nodes[event.Node].Fork].Remove(miningEvent)

		s.Nodes[event.Node].Fork = event.Fork
		s.ScheduleBlockMinedEvent(event.Node, event.Block, event.Depth+1)
		return
	}
	logrus.Debug("BLOCK RECEIVED | ignoring")
}

// func HandleBlockMinedEvent(event *Event) {
// 	// if simulation.Nodes[event.Node].CurrentMiningEvent.Block != event.Block {
// 	// 	return
// 	// }
//
// 	simulation.Blocks[event.Block].FinishedAt = event.DispatchAt
// 	simulation.Blocks[event.Block].Mined = true
//
// 	for currentNode := 0; currentNode < len(simulation.Nodes); currentNode += 1 {
// 		if currentNode == event.Node {
// 			ScheduleBlockMinedEvent(currentNode, event.Block)
// 		} else {
// 			ScheduleBlockReceivedEvent(currentNode, event.Block)
// 		}
// 	}
// }
// func HandleBlockMinedEvent2(event *Event) {
// 	if simulation.Nodes[event.Node].CurrentMiningEvent.Block != event.Block {
// 		return
// 	}
//
// 	simulation.Blocks[event.Block].FinishedAt = event.DispatchAt
// 	simulation.Blocks[event.Block].Mined = true
//
// 	for currentNode := 0; currentNode < len(simulation.Nodes); currentNode += 1 {
// 		if currentNode == event.Node {
// 			ScheduleBlockMinedEvent(currentNode, event.Block)
// 		} else {
//
// 			receivedTimeNext := GetNextReceivedTime()
//
// 			currentMiningEvent := simulation.Nodes[currentNode].CurrentMiningEvent
//
// 			// The received event will happen before the mining event. In this case
// 			// we want to handle everything as it was before, as such are going to
// 			// schedule the Block Received event.
// 			if receivedTimeNext > currentMiningEvent.DispatchAt {
// 				simulation.Queue.Push(&Event{
// 					Type:       BlockReceivedEvent,
// 					Node:       currentNode,
// 					Block:      event.Block,
// 					DispatchAt: receivedTimeNext,
// 				})
// 				continue
// 			}
// 			// otherwise, we are going to have to cancel the original mining event anyway,
// 			// and this is a much more common occurance, that we are not going to schedule
// 			// a new element, for GC to allocate and collect later, which should save us
// 			// in speed
// 			minedBlock := simulation.Nodes[currentNode].CurrentMiningEvent.Block
//
// 			if simulation.Blocks[minedBlock].PreviousBlock == simulation.Blocks[event.Block].PreviousBlock {
// 				// simulation.Queue.Remove(simulation.Nodes[currentNode].CurrentMiningEvent)
// 				simulation.Blocks[minedBlock].FinishedAt = receivedTimeNext
// 				oldTime := simulation.CurrentTime
// 				simulation.CurrentTime = receivedTimeNext
// 				ScheduleBlockMinedEvent(event.Node, event.Block)
// 				simulation.CurrentTime = oldTime
// 				continue
// 			}
//
// 			if simulation.Blocks[minedBlock].Depth+simulation.Configuration.ChainReogranizationThreshold <= simulation.Blocks[event.Block].Depth {
// 				// simulation.Queue.Remove(simulation.Nodes[currentNode].CurrentMiningEvent)
// 				simulation.Blocks[minedBlock].FinishedAt = receivedTimeNext
// 				oldTime := simulation.CurrentTime
// 				simulation.CurrentTime = receivedTimeNext
// 				ScheduleBlockMinedEvent(event.Node, event.Block)
// 				simulation.CurrentTime = oldTime
// 				continue
// 			}
//
// 		}
// 	}
// }
//
