package internal

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"github.com/sirupsen/logrus"
)

func NewSQLite(path string) *SQLite {
	database, err := sql.Open("sqlite3", path)
	if err != nil {
		logrus.Fatal(err)
	}

	if err = database.Ping(); err != nil {
		logrus.Fatal(err)
	}

	sqlite := &SQLite{
		_database: database,
		_path:     path,
	}

	_, _ = database.Exec("PRAGMA journal_mode = WAL;")
	_, _ = database.Exec("PRAGMA synchronous = NORMAL;")
	_, _ = database.Exec("PRAGMA temp_store = MEMORY;")
	_, _ = database.Exec("PRAGMA locking_mode = EXCLUSIVE;")

	sqlite.PrepareBlocksTable()
	sqlite.PrepareNodesTable()
	sqlite.PrepareProofOfWorkConsensusTable()
	sqlite.PrepareSlimcoinProofOfBurnConsensusTable()
	sqlite.PrepareRazerProofOfBurnConsensusTable()
	sqlite.PreparePricingProofOfBurnConsensusTable()
	sqlite.PreparePricingProofOfBurnBurnTransactionTable()

	return sqlite
}

func (s *SQLite) Close() {

	s._blocks.Flush()
	s._nodes.Flush()
	s._proofOfWorkConsensus.Flush()
	s._slimcoinProofOfBurnConsensus.Flush()
	s._razerProofOfBurnConsensus.Flush()
	s._pricingProofOfBurnConsensus.Flush()
	s._pricingProofOfBurnBurnTransactions.Flush()
	s._database.Close()
}

const batchSize = 500

type SQLite struct {
	_database *sql.DB
	_path     string

	_blocks                             *SQLiteTable[BlocksTable]
	_nodes                              *SQLiteTable[NodesTable]
	_proofOfWorkConsensus               *SQLiteTable[ProofOfWorkConsensusTable]
	_slimcoinProofOfBurnConsensus       *SQLiteTable[SlimcoinProofOfBurnConsensusTable]
	_razerProofOfBurnConsensus          *SQLiteTable[RazerProofOfBurnConsensusTable]
	_pricingProofOfBurnConsensus        *SQLiteTable[PricingProofOfBurnBurnConsensusTable]
	_pricingProofOfBurnBurnTransactions *SQLiteTable[PricingProofOfBurnBurnTransactionTable]
}
