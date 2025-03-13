package internal

type Block struct {
	Node          *Node
	PreviousBlock *Block
	MinedAt       float64
	Depth         int
}

type Node struct {
	CurrentlyMinedBlock *Block
}

type EventType int

const (
	BlockMinedEvent    EventType = 0
	BlockReceivedEvent EventType = 1
)

type Event struct {
	Type EventType

	Block *Block
	Node  *Node

	DispatchAt float64

	Index int
}
