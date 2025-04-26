package internal

type ConsensusType int64

const (
	Genesis     ConsensusType = -1
	ProofOfWork ConsensusType = 0
	ProofOfBurn ConsensusType = 1
)

type Consensus interface {
	Initialize()

	Adjust(event *Event_BlockMined)

	GetType() ConsensusType

	CanMine(receivedEvent *Event_BlockReceived) bool
	GetNextMiningTime(event *Event_BlockMined) float64

	Synchronize(consensus Consensus)
}
