package internal

type Block struct {
	Node          int
	PreviousBlock int
	MinedAt       float64
	Depth         int
}

type Node struct {
	CurrentlyMinedBlock int
}

type EventType int

const (
	BlockMinedEvent    EventType = 0
	BlockReceivedEvent EventType = 1
)

type Event struct {
	Type EventType

	Block int
	Node  int

	DispatchAt float64

	Index int
}
