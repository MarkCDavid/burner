package internal

type Node struct {
	CurrentEvent *Event
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
