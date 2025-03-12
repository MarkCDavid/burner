package simulator

import (
	"github.com/MarkCDavid/burner/internal/pq"
	randomness "github.com/MarkCDavid/burner/internal/random"
)

type Simulator struct {
  EventQueue pq.PriorityQueue
  Randomness *randomness.Randomness
}
