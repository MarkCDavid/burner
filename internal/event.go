package internal

import (
	"encoding/json"
)

type EventType int

const (
	BlockMinedEvent    EventType = 0
	BlockReceivedEvent EventType = 1
)

type BlockType int

const (
	Genesis             BlockType = -1
	ProofOfWork         BlockType = 0
	SlimcoinProofOfBurn BlockType = 1
)

type Block struct {
	Id           int64
	Node         int64
	Depth        int64
	Transactions int64
	Type         BlockType
}

type Event struct {
	Node int64

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
	return e.DispatchAt - e.ScheduledAt
}
