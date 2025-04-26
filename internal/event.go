package internal

//
// import (
// 	"encoding/json"
// )

type Event interface {
	GetIndex() int
	SetIndex(index int)
	EventTime() float64
	Duration() float64
	Handle()
}

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
