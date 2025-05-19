package internal

import (
	"fmt"
)

type BlocksTable struct {
	Id                 int64
	PreviousBlockId    *int64
	MinedBy            int64
	Depth              int64
	StartedAt          float64
	FinishedAt         float64
	PreviousFinishedAt float64
	Abandoned          bool
	Transactions       int64
	BlockType          ConsensusType
}

func (db *SQLite) PrepareBlocksTable() {
	tableName := "blocks"
	createTable := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id INTEGER PRIMARY KEY, 
		previousBlockId INTEGER,
		minedBy INTEGER NOT NULL,
		depth INTEGER NOT NULL,
		startedAt REAL NOT NULL,
		finishedAt REAL NOT NULL,
		previousFinishedAt REAL NOT NULL,
		abandoned BOOLEAN NOT NULL,
		transactions INTEGER NOT NULL,
		blockType INTEGER NOT NULL,
		FOREIGN KEY (previousBlockId) REFERENCES blocks(id)
	);`, tableName,
	)

	insertInto := fmt.Sprintf(`
	INSERT INTO %s (
    id,
		previousBlockId,
		minedBy,
		depth,
		startedAt,
		finishedAt,
		previousFinishedAt,
		abandoned,
		transactions,
		blockType
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`,
		tableName,
	)

	db._blocks = NewSQLiteTable[BlocksTable](
		db,
		createTable,
		insertInto,
	)
}

func (db *SQLite) SaveBlock(event *Event_BlockMined) {
	var previousBlockId *int64
	if event.PreviousBlock != nil {
		previousBlockId = &event.PreviousBlock.Id
	}
	db._blocks.Save(BlocksTable{
		Id:                 event.Block.Id,
		PreviousBlockId:    previousBlockId,
		MinedBy:            event.Block.Node.Id,
		Depth:              event.Block.Depth,
		StartedAt:          event.Block.StartedAt,
		FinishedAt:         event.Block.FinishedAt,
		PreviousFinishedAt: event.PreviousBlock.FinishedAt,
		Abandoned:          event.Block.Abandoned,
		Transactions:       event.Block.Transactions,
		BlockType:          event.Block.Consensus.GetType(),
	})

}
