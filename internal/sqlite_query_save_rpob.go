package internal

import (
	"fmt"
)

type RazerProofOfBurnConsensusTable struct {
	SimulationTime float64
	NodeId         int64
	Chance         float64
	UpdateType     ConsensusUpdateType
}

func (db *SQLite) PrepareRazerProofOfBurnConsensusTable() {
	tableName := "razer_proof_of_burn_consensus"
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

	db._razerProofOfBurnConsensus = NewSQLiteTable[RazerProofOfBurnConsensusTable](
		db,
		createTable,
		insertInto,
	)
}

func (db *SQLite) SaveRazerProofOfBurnConsensus(consensus Consensus, updateType ConsensusUpdateType) {
	rpobConsensus, ok := consensus.(*Consensus_RPoB)
	if !ok {
		return
	}

	db._slimcoinProofOfBurnConsensus.Save(SlimcoinProofOfBurnConsensusTable{
		SimulationTime: rpobConsensus.Node.Simulation.CurrentTime,
		NodeId:         rpobConsensus.Node.Id,
		Chance:         float64(1) / rpobConsensus.Difficulty,
		UpdateType:     updateType,
	})

}
