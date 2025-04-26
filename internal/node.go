package internal

func NewNode(simulation *Simulation) *Node {

	node := &Node{
		Id:         int64(len(simulation.Nodes)),
		Simulation: simulation,
	}

	AddConsensus_PPoB(node)
	AddConsensus_RPoB(node)
	AddConsensus_SPoB(node)
	AddConsensus_PoW(node)

	return node
}

func (to *Node) SynchronizeConsensus(from *Node) {
	if to.ProofOfWork != nil {
		to.ProofOfWork.Synchronize(from.ProofOfWork)
	}

	if to.ProofOfBurn != nil {
		to.ProofOfBurn.Synchronize(from.ProofOfBurn)
	}
}

type Node struct {
	Id int64

	Simulation *Simulation

	Event         *Event_BlockMined
	PreviousBlock *Block

	ProofOfWork Consensus
	ProofOfBurn Consensus

	Transactions int64
}

func (n *Node) GetConsensus(event *Event_BlockReceived) Consensus {
	if n.ProofOfWork != nil && n.ProofOfWork.CanMine(event) {
		return n.ProofOfWork
	}

	if n.ProofOfBurn != nil && n.ProofOfBurn.CanMine(event) {
		return n.ProofOfBurn
	}

	return nil
}
