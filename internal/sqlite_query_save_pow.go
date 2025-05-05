package internal

import (
	"fmt"
)

type ProofOfWorkConsensusTable struct {
	SimulationTime float64
	NodeId         int64
	Difficulty     float64
	UpdateType     ConsensusUpdateType
}

func (db *SQLite) PrepareProofOfWorkConsensusTable() {
	tableName := "proof_of_work_consensus"
	createTable := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp REAL NOT NULL,
    nodeId INTEGER NOT NULL,
    difficulty REAL NOT NULL,
    eventType INTEGER NOT NULL
  );`, tableName,
	)

	insertInto := fmt.Sprintf(`INSERT INTO %s (
    timestamp,
    nodeId,
    difficulty,
    eventType
  ) VALUES (?, ?, ?, ?);`,
		tableName,
	)

	db._proofOfWorkConsensus = NewSQLiteTable[ProofOfWorkConsensusTable](
		db,
		createTable,
		insertInto,
	)
}

func (db *SQLite) SaveProofOfWorkConsensus(consensus Consensus, updateType ConsensusUpdateType) {
	powConsensus, ok := consensus.(*Consensus_PoW)
	if !ok {
		return
	}

	db._proofOfWorkConsensus.Save(ProofOfWorkConsensusTable{
		SimulationTime: powConsensus.Node.Simulation.CurrentTime,
		NodeId:         powConsensus.Node.Id,
		Difficulty:     powConsensus.Difficulty,
		UpdateType:     updateType,
	})

}
