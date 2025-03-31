package internal

import (
	"fmt"
	"os"
)

func ExportBlocksToDotGraph(blocks []Block, onlyMined bool, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintln(file, "digraph G {")

	for i := 0; i < len(blocks); i++ {
		if !onlyMined {
			addNode(file, i)
		} else if i == 0 || blocks[i].Mined {
			addNode(file, i)
		}
	}

	for i := 0; i < len(blocks); i++ {
		if !onlyMined {
			addEdge(file, i)
		} else if i == 0 || blocks[i].Mined {
			addEdge(file, i)
		}
	}

	fmt.Fprintln(file, "}")

	return nil
}

func addNode(file *os.File, blockId int) {
	color := "white"
	if blockId == 0 || pobSimulation.Blocks[blockId].Mined {
		switch pobSimulation.Blocks[blockId].BlockType {
		case ProofOfBurn:
			color = "red"
		case ProofOfWork:
			color = "green"
		default:
			panic("unknown block type")
		}
	}
	miningTime := 0.0
	if pobSimulation.Blocks[blockId].FinishedAt > 0 {
		miningTime = (pobSimulation.Blocks[blockId].FinishedAt - pobSimulation.Blocks[blockId].StartedAt)

	}
	fmt.Fprintf(
		file,
		"    \"%d\" [label=\"block: %d\nmined by: %d\nmined for: %f\ndiff: %f\" style=\"filled\" fillcolor=\"%s\" constraint=\"false\"];\n",
		blockId,
		blockId,
		pobSimulation.Blocks[blockId].Node,
		miningTime,
		pobSimulation.Blocks[blockId].ProofOfBurnDifficulty,
		color,
	)
}

func addEdge(file *os.File, blockId int) {
	if pobSimulation.Blocks[blockId].PreviousBlock == -1 {
		return
	}

	fmt.Fprintf(
		file,
		"    \"%d\" -> \"%d\";\n",
		pobSimulation.Blocks[blockId].PreviousBlock,
		blockId,
	)
}
