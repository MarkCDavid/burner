package internal

type ConsensusType int64

const (
	Genesis     ConsensusType = -1
	ProofOfWork ConsensusType = 0
	ProofOfBurn ConsensusType = 1
)

type Consensus interface {
	Initialize()

	Adjust(event *Event)

	GetType() ConsensusType

	CanMine(receivedEvent *Event) bool
	GetNextMiningTime(event *Event) float64

	Synchronize(consensus Consensus)
}
