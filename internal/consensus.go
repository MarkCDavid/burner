package internal

type ConsensusType int64

const (
	Genesis     ConsensusType = -1
	ProofOfWork ConsensusType = 0
	ProofOfBurn ConsensusType = 1
)

type ConsensusUpdateType int64

const (
	Initialize  ConsensusUpdateType = -1
	Adjust      ConsensusUpdateType = 0
	Synchronize ConsensusUpdateType = 1
)

func (t ConsensusType) ToString() string {
	switch t {
	case Genesis:
		return "Genesis"
	case ProofOfWork:
		return "ProofOfWork"
	case ProofOfBurn:
		return "ProofOfBurn"
	}
	return "N/A"
}

type Consensus interface {
	Initialize()

	GetType() ConsensusType

	CanMine(event Event) bool
	Adjust(event Event)

	GetNextMiningTime(event *Event_BlockMined) float64

	Synchronize(consensus Consensus)
}
