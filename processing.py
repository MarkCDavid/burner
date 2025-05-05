import os
import duckdb
import sys


def open_database(path):
    name = os.path.splitext(os.path.basename(path))[0]
    db = duckdb.connect()
    db.execute("INSTALL sqlite;")
    db.execute("LOAD sqlite;")
    db.execute(f"ATTACH '{path}' (TYPE sqlite);")
    db.execute(f"USE '{name}';")

    return db


QUERY_BLOCKS_PER_NODE = """
    SELECT 
        minedBy, 
        blockType,
        COUNT(*) as block_count
    FROM blocks
    GROUP BY blockType, minedBy
    ORDER BY block_count DESC
"""

QUERY_BLOCK_PRODUCTION_TIMES = """
    SELECT 
        blockType,
        COUNT(*) AS count,
        AVG(finishedAt - startedAt) AS avg_block_time,
        MIN(finishedAt - startedAt) AS min_block_time,
        MAX(finishedAt - startedAt) AS max_block_time
    FROM blocks
    GROUP BY blockType
    ORDER BY avg_block_time DESC;
"""

QUERY_GET_LEAF_BLOCKS = """
WITH RECURSIVE valid_chain(id, previousBlockId) AS (
  SELECT id, previousBlockId
  FROM blocks
  WHERE previousBlockId IS NULL AND abandoned = 0  -- Start from genesis

  UNION ALL

  SELECT b.id, b.previousBlockId
  FROM blocks b
  JOIN valid_chain vc ON b.previousBlockId = vc.id
  WHERE b.abandoned = 0
)
SELECT * FROM valid_chain;
"""


if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python script.py <path_to_sqlite_file>")
        sys.exit(1)

    db = open_database(sys.argv[1])

    df = db.execute(QUERY_GET_LEAF_BLOCKS).fetchdf()
    print(df)
