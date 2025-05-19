package internal

import (
	"fmt"
)

type LabelTable struct {
	Label string
}

func (db *SQLite) PrepareLabelTable() {
	tableName := "label"
	createTable := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id INTEGER PRIMARY KEY, 
		label TEXT 
	);`, tableName)

	insertInto := fmt.Sprintf(`
	INSERT INTO %s (
		label
	) VALUES (?);`, tableName)

	db._label = NewSQLiteTable[LabelTable](
		db,
		createTable,
		insertInto,
	)
}

func (db *SQLite) SaveLabel(label string) {
	db._label.Save(LabelTable{
		Label: label,
	})
}
