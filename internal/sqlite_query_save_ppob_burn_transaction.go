package internal

import (
	"fmt"
)

type PricingProofOfBurnBurnTransactionTable struct {
	NodeId    int64
	BurnedAt  int64
	BurnedFor float64
}

func (db *SQLite) PreparePricingProofOfBurnBurnTransactionTable() {
	tableName := "pricing_proof_of_burn_burn_transaction"
	createTable := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nodeId INTEGER NOT NULL,
    burnedAt INTEGER NOT NULL,
    burnedFor REAL NOT NULL
  );`, tableName,
	)

	insertInto := fmt.Sprintf(`INSERT INTO %s (
		nodeId,
		burnedAt,
		burnedFor
  ) VALUES (?, ?, ?);`,
		tableName,
	)

	db._pricingProofOfBurnBurnTransactions = NewSQLiteTable[PricingProofOfBurnBurnTransactionTable](
		db,
		createTable,
		insertInto,
	)
}

func (db *SQLite) SavePricingProofOfBurnBurnTransaction(burnTransaction BurnTransaction) {

	db._pricingProofOfBurnBurnTransactions.Save(PricingProofOfBurnBurnTransactionTable{
		NodeId:    burnTransaction.BurnedBy.Id,
		BurnedAt:  burnTransaction.BurnedAt,
		BurnedFor: burnTransaction.BurnedFor,
	})
}
