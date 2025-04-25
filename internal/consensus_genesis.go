package internal

type Consensus_Genesis struct{}

func (*Consensus_Genesis) Initialize() {}

func (*Consensus_Genesis) Adjust(event *Event) {}

func (*Consensus_Genesis) GetType() ConsensusType { return Genesis }

func (*Consensus_Genesis) CanMine(receivedEvent *Event) bool      { return false }
func (*Consensus_Genesis) GetNextMiningTime(event *Event) float64 { return 0 }

func (*Consensus_Genesis) Synchronize(consensus Consensus) {}
