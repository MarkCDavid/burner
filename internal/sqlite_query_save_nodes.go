package internal

import (
	"fmt"
)

type NodesTable struct {
	Id    int64
	Power float64
}

func (db *SQLite) PrepareNodesTable() {
	tableName := "nodes"
	createTable := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id INTEGER PRIMARY KEY, 
		power REAL NOT NULL
	);`, tableName,
	)

	insertInto := fmt.Sprintf(`
	INSERT INTO %s (
    id,
		power
	) VALUES (?, ?);`,
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
		Id:    node.Id,
		Power: node.Power,
	})

}
