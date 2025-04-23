package internal

import "encoding/json"

type Node struct {
	CurrentEvent *Event
}

type EventType int

const (
	BlockMinedEvent    EventType = 0
	BlockReceivedEvent EventType = 1
)

type Event struct {
	Node      int       `json:"node"`
	EventType EventType `json:"eventType"`

	Block         int `json:"block"`
	PreviousBlock int `json:"previousBlock"`
	Depth         int `json:"depth"`

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
