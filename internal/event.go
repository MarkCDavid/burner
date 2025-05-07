package internal

type Event interface {
	GetIndex() int
	SetIndex(index int)
	EventTime() float64
	Handle()
}

type Block struct {
	Id int64

	Node *Node

	Depth        int64
	Transactions int64

	StartedAt  float64
	FinishedAt float64

	Abandoned bool

	Consensus Consensus
}
