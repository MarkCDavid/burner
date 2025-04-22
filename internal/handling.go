package internal

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

func HandleBlockMinedEvent(event *Event) {
	if simulation.Nodes[event.Node].CurrentEvent.Block != event.Block {
		return
	}

	logrus.Debug("================= Handle Block Mined Event")
	ScheduleBlockMinedEvent(event.Node, event.Block)
	for currentNode := 0; currentNode < len(simulation.Nodes); currentNode += 1 {
		if currentNode != event.Node {

			logrus.Debugf("====== Scheduling for %d by %d", currentNode, event.Node)
			ScheduleBlockReceivedEvent(currentNode, event)
		}
	}
}

func HandleBlockReceivedEvent(event *Event) {
	miningEvent := simulation.Nodes[event.Node].CurrentEvent

	eventJson, _ := json.Marshal(event)
	miningEventJson, _ := json.Marshal(miningEvent)
	logrus.Debugf("Received event: %s", string(eventJson))
	logrus.Debugf("Mining event: %s", string(miningEventJson))

	if event.PreviousBlock == miningEvent.PreviousBlock {
		logrus.Debug("HandleBlockReceivedEvent: Same Block")
		logrus.Debug(miningEvent)
		// Stats will be calculated when scheduling new block.
		ScheduleBlockMinedEvent(event.Node, event.Block)
		return
	}
	if miningEvent.Depth+simulation.Configuration.ChainReogranizationThreshold <= event.Depth {
		logrus.Debug("HandleBlockReceivedEvent: Enough Depth")
		logrus.Debug(miningEvent)
		// Stats will be calculated when scheduling new block.
		ScheduleBlockMinedEvent(event.Node, event.Block)
		return
	}
	logrus.Debug("HandleBlockReceivedEvent: Ignoring")
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
