package internal

// import (
// 	"fmt"
// 	"os"
//
// 	"github.com/sirupsen/logrus"
// )
//
// func (s *Simulation) ExportBlocksToDotGraph(filename string) error {
// 	file, err := os.Create(filename)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()
//
// 	fmt.Fprintln(file, "digraph G {")
//
// 	fmt.Fprintf(file, "    \"seed[%d]\"", s.Random.GetSeed())
//
// 	for i := range s.Blocks {
// 		if s.Blocks[i].Mined {
// 			logrus.Debug(s.Blocks[i].ToString())
// 			s.addNode(file, i)
// 		}
// 	}
//
// 	for i := range s.Blocks {
// 		if s.Blocks[i].Mined {
// 			s.addEdge(file, i)
// 		}
// 	}
//
// 	fmt.Fprintln(file, "}")
//
// 	return nil
// }
//
// func (s *Simulation) addNode(file *os.File, blockId int) {
// 	color := "white"
// 	if blockId == 0 || s.Blocks[blockId].Mined {
// 		color = "green"
// 	}
//
// 	fmt.Fprintf(
// 		file,
// 		"    \"%d\" [label=\"block: %d\nmined by: %d\nfork: %d\nscheduled at: %.2f\ndispatch [%.2f | %.2f]\" style=\"filled\" fillcolor=\"%s\"];\n",
// 		s.Blocks[blockId].Block,
// 		s.Blocks[blockId].Block,
// 		s.Blocks[blockId].Node,
// 		s.Blocks[blockId].Fork,
// 		s.Blocks[blockId].ScheduledAt,
// 		s.Blocks[blockId].DispatchAt,
// 		s.Blocks[blockId].DispatchAt+s.Configuration.UpperBoundNetworkLatency,
// 		color,
// 	)
// }
//
// func (s *Simulation) addEdge(file *os.File, blockId int) {
// 	if s.Blocks[blockId].PreviousBlock == -1 {
// 		return
// 	}
//
// 	fmt.Fprintf(
// 		file,
// 		"    \"%d\" -> \"%d\";\n",
// 		s.Blocks[blockId].PreviousBlock,
// 		s.Blocks[blockId].Block,
// 	)
// }
