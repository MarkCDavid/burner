package internal

import (
	"fmt"
)

type PricingProofOfBurnBurnConsensusTable struct {
	SimulationTime float64
	NodeId         int64
	Price          float64
}

func (db *SQLite) PreparePricingProofOfBurnConsensusTable() {
	tableName := "pricing_proof_of_burn_burn_consensus"
	createTable := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp REAL NOT NULL,
    nodeId INTEGER NOT NULL,
    price REAL NOT NULL
  );`, tableName,
	)

	insertInto := fmt.Sprintf(`INSERT INTO %s (
		timestamp,
		nodeId,
		price
  ) VALUES (?, ?, ?);`,
		tableName,
	)

	db._pricingProofOfBurnConsensus = NewSQLiteTable[PricingProofOfBurnBurnConsensusTable](
		db,
		createTable,
		insertInto,
	)
}

func (db *SQLite) SavePricingProofOfBurnConsensus(consensus Consensus, updateType ConsensusUpdateType) {
	ppobConsensus, ok := consensus.(*Consensus_PPoB)
	if !ok {
		return
	}
	db._pricingProofOfBurnConsensus.Save(PricingProofOfBurnBurnConsensusTable{
		SimulationTime: ppobConsensus.Node.Simulation.CurrentTime,
		NodeId:         ppobConsensus.Node.Id,
		Price:          ppobConsensus.Price,
	})
}
