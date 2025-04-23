package internal

import "encoding/json"

type Node struct {
	CurrentEvent *Event
	Fork         int
}

type EventType int

const (
	BlockMinedEvent    EventType = 0
	BlockReceivedEvent EventType = 1
)

type Block struct {
	Node int `json:"node"`

	Block         int  `json:"block"`
	PreviousBlock int  `json:"previousBlock"`
	Depth         int  `json:"depth"`
	Fork          int  `json:"fork"`
	Mined         bool `json:"mined"`

	ScheduledAt float64 `json:"scheduledAt"`
	DispatchAt  float64 `json:"dispatchAt"`
}

type Event struct {
	Node int `json:"node"`

	Block         int `json:"block"`
	PreviousBlock int `json:"previousBlock"`
	Depth         int `json:"depth"`
	Fork          int `json:"fork"`

	ScheduledAt float64 `json:"scheduledAt"`
	DispatchAt  float64 `json:"dispatchAt"`

	Index int `json:"index"`
}

func (e *Event) ToString() string {
	if e == nil {
		return "nil"
	}
	result, _ := json.Marshal(*e)
	return string(result)
}

func (b *Block) ToString() string {
	if b == nil {
		return "nil"
	}
	result, _ := json.Marshal(*b)
	return string(result)
}
