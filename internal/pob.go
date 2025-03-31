package internal

type PoBSimulation struct {
	Nodes       []PoBNode
	Blocks      []Block
	Queue       EventQueue
	Random      *Rng
	CurrentTime float64
}

var pobSimulation PoBSimulation

func PoB(seed *int64) error {
	random := CreateRandom(seed)
	nodes := make([]PoBNode, 100)
	for i := 0; i < len(nodes); i++ {
		for j := 0; j < 10; j++ {
			if j == 0 {
				nodes[i].Burners[j] = 0
			} else {
				nodes[i].Burners[j] = -1
			}
		}
	}

	pobSimulation = PoBSimulation{
		Nodes: nodes,
		Blocks: []Block{
			{
				Node:                  -1,
				PreviousBlock:         -1,
				Depth:                 0,
				ProofOfBurnDifficulty: 1.0,
			},
		},
		Queue:       CreateEventQueue(),
		Random:      random,
		CurrentTime: 0,
	}

	for nodeId := 0; nodeId < 100; nodeId += 1 {
		ScheduleBlockMinedEventPoB(nodeId, 0)
	}

	for pobSimulation.Queue.Len() > 0 {
		event := pobSimulation.Queue.Pop()
		pobSimulation.CurrentTime = event.DispatchAt

		switch event.Type {
		case BlockMinedEvent:
			HandleBlockMinedEventPoB(event)
		case BlockReceivedEvent:
			HandleBlockReceivedEventPoB(event)
		default:
			panic("Unknown event")
		}

		if pobSimulation.CurrentTime > 1000000 {
			break
		}
	}
	ExportBlocksToDotGraph(pobSimulation.Blocks, true, "chain.dot")

	return nil
}
