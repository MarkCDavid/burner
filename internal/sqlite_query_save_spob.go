package internal

import (
	"fmt"
)

type SlimcoinProofOfBurnConsensusTable struct {
	SimulationTime float64
	NodeId         int64
	Chance         float64
	UpdateType     ConsensusUpdateType
}

func (db *SQLite) PrepareSlimcoinProofOfBurnConsensusTable() {
	tableName := "slimcoin_proof_of_burn_consensus"
	createTable := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp REAL NOT NULL,
    nodeId INTEGER NOT NULL,
    chance REAL NOT NULL,
    eventType INTEGER NOT NULL
  );`, tableName,
	)

	insertInto := fmt.Sprintf(`INSERT INTO %s (
    timestamp,
    nodeId,
    chance,
    eventType
  ) VALUES (?, ?, ?, ?);`,
		tableName,
	)

	db._slimcoinProofOfBurnConsensus = NewSQLiteTable[SlimcoinProofOfBurnConsensusTable](
		db,
		createTable,
		insertInto,
	)
}

func (db *SQLite) SaveSlimcoinProofOfBurnConsensus(consensus Consensus, updateType ConsensusUpdateType) {
	spobConsensus, ok := consensus.(*Consensus_SPoB)
	if !ok {
		return
	}

	db._slimcoinProofOfBurnConsensus.Save(SlimcoinProofOfBurnConsensusTable{
		SimulationTime: spobConsensus.Node.Simulation.CurrentTime,
		NodeId:         spobConsensus.Node.Id,
		Chance:         float64(1) / spobConsensus.Difficulty,
		UpdateType:     updateType,
	})

}
