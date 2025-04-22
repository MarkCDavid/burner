package internal

import (
	"github.com/sirupsen/logrus"
)

func GetNextMiningTime() float64 {
	return simulation.CurrentTime + simulation.Random.Expovariate(simulation.Configuration.InverseAverageBlockFrequency)
}

func GetNextReceivedTime() float64 {
	if simulation.Configuration.AverageNetworkLatencyInSeconds <= 0 {
		return simulation.CurrentTime
	}

	delta := simulation.Random.Expovariate(simulation.Configuration.InverseAverageNetworkLatency)
	if delta > simulation.Configuration.UpperBoundNetworkLatency {
		delta = simulation.Configuration.UpperBoundNetworkLatency
	}

	if delta < simulation.Configuration.LowerBoundNetworkLatency {
		delta = simulation.Configuration.LowerBoundNetworkLatency
	}

	return simulation.CurrentTime + delta
}

func ScheduleBlockMinedEvent(
	minedBy int,
	previousBlock int,
) {

	logrus.Debugf("============ Schedule Block Mined Event (%d)", minedBy)
	simulation.BlockCount += 1
	minedAt := GetNextMiningTime()

	nextBlockMinedEvent := simulation.BlockMinedEventQueue.Peek()

	if simulation.BlockMinedEventQueue.Len() == 1 && nextBlockMinedEvent != nil {
		newBlockLowerBound := minedAt
		newBlockUpperBound := minedAt + simulation.Configuration.UpperBoundNetworkLatency

		oldBlockLowerBound := nextBlockMinedEvent.DispatchAt
		oldBlockUpperBound := nextBlockMinedEvent.DispatchAt + simulation.Configuration.UpperBoundNetworkLatency

		// The newly scheduled block here might be:
		if newBlockUpperBound < oldBlockLowerBound {

			logrus.Infof("Block mined earlier than earlierst receival time of most recent block. Remove old, calculate statistics, schedule new.")
			// 1. Scheduled before the next block, remove old, stats for old
			simulation.BlockMinedEventQueue.Remove(nextBlockMinedEvent)
			// fmt.Println("too early, remove later, just calculate stats")
			// Calculate stats here?
			scheduleBlockMinedEvent(minedBy, minedAt, previousBlock)
		} else if newBlockLowerBound > oldBlockUpperBound {
			// 2. Scheduled after the next block, don't schedule, just stats
			logrus.Infof("Block mined later than latest receival time of most recent block. Calculate statistics, skip.")
			// IMPORTANT: Need to calculate receive delay here - across many nodes it matters.
		} else {

			logrus.Infof("Blocks mined within block receival interval.")
			// 3. Scheduled during the next block
			scheduleBlockMinedEvent(minedBy, minedAt, previousBlock)
			// panic("wowza, we gon have to figure out forks")
		}
	} else {

		logrus.Infof("Regular scheduling")
		scheduleBlockMinedEvent(minedBy, minedAt, previousBlock)
	}
}

func scheduleBlockMinedEvent(minedBy int, minedAt float64, previousBlock int) {
	previousEvent := simulation.Nodes[minedBy].CurrentEvent

	depth := 0
	if previousEvent != nil {
		depth = previousEvent.Depth + 1
	}

	event := &Event{
		Node: minedBy,

		Block:         simulation.BlockCount,
		PreviousBlock: previousBlock,
		Depth:         depth,

		ScheduledAt: simulation.CurrentTime,
		DispatchAt:  minedAt,
	}

	simulation.Nodes[minedBy].CurrentEvent = event
	simulation.BlockMinedEventQueue.Push(event)
}

func ScheduleBlockReceivedEvent(receivedBy int, minedEvent *Event) {
	// If only one event in the BlockReceivedEventQueue - don't actually handle block received events.
	if simulation.BlockMinedEventQueue.Len() > 1 {
		simulation.BlockReceivedEventQueue.Push(&Event{
			Node: receivedBy,

			Block:         minedEvent.Block,
			PreviousBlock: minedEvent.PreviousBlock,
			Depth:         minedEvent.Depth,

			ScheduledAt: simulation.CurrentTime,
			DispatchAt:  GetNextReceivedTime(),
		})

	} else {
		logrus.Info("doing out of band calculation")
		// Here we would need to calculate the statistics for semi mined block
		// that was dismissed before too.

		// calculate when the block would be received

		oldSimulationTime := simulation.CurrentTime
		simulation.CurrentTime = GetNextReceivedTime()
		ScheduleBlockMinedEvent(receivedBy, minedEvent.Block)
		simulation.CurrentTime = oldSimulationTime
	}

}
