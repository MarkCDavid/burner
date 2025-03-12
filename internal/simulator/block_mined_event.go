package simulator

import (
	"os"

	"github.com/MarkCDavid/burner/internal/model"
	"github.com/MarkCDavid/burner/internal/simulator"
)

type BlockMinedEvent struct {
  Miner *model.Node
  Block *model.Block
}

func (event *BlockMinedEvent) Handle(simulator Simulator) {
  lastBlock := event.Miner.Blockchain[len(event.Miner.Blockchain) - 1]
  if lastBlock != event.Block.PreviousBlock {
    return
  }

  event.Miner.Blockchain = append(event.Miner.Blockchain, event.Block)

  // Schedule block receival events
}

func ScheduleBlockMinedEvent(
  simulator Simulator, 
  miner *model.Node,
  previousBlock *model.Block,
) {
  event := &BlockMinedEvent {
    Miner: miner,
    Block: &model.Block{
      BlockId: simulator.Randomness.Id(),
      PreviousBlock: previousBlock,
      Miner: miner,
      // Fix hardcoded values
      MinedAt: previousBlock.MinedAt + simulator.Randomness.Expovariate(1.0 / 10.0),
    },
  }

}
