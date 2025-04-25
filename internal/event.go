package internal

import (
	"encoding/json"
)

type EventType int

const (
	BlockMinedEvent    EventType = 0
	BlockReceivedEvent EventType = 1
)

type Block struct {
	Id           int64
	Node         *Node
	Depth        int64
	Transactions int64
	Consensus    Consensus
}

type Event struct {
	Node *Node

	EventType EventType

	Block         *Block
	PreviousBlock *Block

	ScheduledAt float64
	DispatchAt  float64

	Index int
}

func (e *Event) ToString() string {
	if e == nil {
		return "nil"
	}
	result, _ := json.Marshal(*e)
	return string(result)
}

func (e *Event) Duration() float64 {
	return e.Node.Simulation.CurrentTime - e.ScheduledAt
}

func (e *Event) PowerUsed() float64 {
	return e.Duration() * e.Node.Power[e.Block.Consensus.GetType()]
}

func (e *Event) SetMiner(n *Node) {
	e.Node = n
	n.Event = e
}

func (e *Event) SetReceiver(n *Node) {
	e.Node = n
}
