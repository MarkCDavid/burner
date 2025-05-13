package internal

import (
	"fmt"
)

type NodesTable struct {
	Id           int64
	PowerFull    float64
	PowerIdle    float64
	Transactions int64
}

func (db *SQLite) PrepareNodesTable() {
	tableName := "nodes"
	createTable := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id INTEGER PRIMARY KEY, 
		powerFull REAL NOT NULL,
		powerIdle REAL NOT NULL,
		transactions INTEGER NOT NULL
	);`, tableName,
	)

	insertInto := fmt.Sprintf(`
	INSERT INTO %s (
    id,
		powerFull,
		powerIdle,
		transactions
	) VALUES (?, ?, ?, ?);`,
		tableName,
	)

	db._nodes = NewSQLiteTable[NodesTable](
		db,
		createTable,
		insertInto,
	)
}

func (db *SQLite) SaveNode(node *Node) {
	db._nodes.Save(NodesTable{
		Id:           node.Id,
		PowerFull:    node.PowerFull,
		PowerIdle:    node.PowerIdle,
		Transactions: node.Transactions,
	})

}
