package internal

import (
	"fmt"
)

type NodesTable struct {
	Id               int64
	ProofOfWorkPower float64
	ProofOfBurnPower float64
}

func (db *SQLite) PrepareNodesTable() {
	tableName := "nodes"
	createTable := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id INTEGER PRIMARY KEY, 
		powPower REAL NOT NULL,
		pobPower REAL NOT NULL
	);`, tableName,
	)

	insertInto := fmt.Sprintf(`
	INSERT INTO %s (
    id,
		powPower,
		pobPower
	) VALUES (?, ?, ?);`,
		tableName,
	)

	db._nodes = NewSQLiteTable[NodesTable](
		db,
		createTable,
		insertInto,
	)
}

func (db *SQLite) SaveNode(node *Node) {
	var powPower float64 = 0
	var pobPower float64 = 0

	if node.ProofOfWork != nil {
		powPower = node.ProofOfWork.GetPower()
	}

	if node.ProofOfBurn != nil {
		pobPower = node.ProofOfBurn.GetPower()
	}
	db._nodes.Save(NodesTable{
		Id:               node.Id,
		ProofOfWorkPower: powPower,
		ProofOfBurnPower: pobPower,
	})

}
