package internal

import unsafe "unsafe"

type BlockType int

const (
	ProofOfWork BlockType = 0
	ProofOfBurn BlockType = 1
)

type Block struct {
	Node                  int
	BlockType             BlockType
	PreviousBlock         int
	Depth                 int
	StartedAt             float64
	FinishedAt            float64
	ProofOfBurnDifficulty float64
	Mined                 bool
}

var _block Block

const BlockSize = int(unsafe.Sizeof(_block))

type Node struct {
	CurrentlyMinedBlock int
	NodePower           float64
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

var _event Event

const EventSize = int(unsafe.Sizeof(_event))
