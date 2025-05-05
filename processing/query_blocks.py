from model import Block
from typing import Generator
from database import Database

QUERY_TRACE_CHAIN = """
WITH RECURSIVE chain AS (
  SELECT b0.id, b0.previousBlockId, b0.minedBy, b0.depth, b0.startedAt, b0.finishedAt, b0.abandoned, b0.transactions, b0.blockType
  FROM blocks b0
  WHERE b0.id = ? AND b0.abandoned = 0

  UNION ALL

  SELECT b1.id, b1.previousBlockId, b1.minedBy, b1.depth, b1.startedAt, b1.finishedAt, b1.abandoned, b1.transactions, b1.blockType
  FROM blocks b1
  JOIN chain c0 ON b1.id = c0.previousBlockId
  WHERE b1.abandoned = 0
)
SELECT * FROM chain;
"""

QUERY_GET_BLOCKS = """
SELECT b0.id, b0.previousBlockId, b0.minedBy, b0.depth, b0.startedAt, b0.finishedAt, b0.abandoned, b0.transactions, b0.blockType
FROM blocks b0
WHERE b0.abandoned = 0
"""

QUERY_GET_ALL_BLOCKS = """
SELECT b0.id, b0.previousBlockId, b0.minedBy, b0.depth, b0.startedAt, b0.finishedAt, b0.abandoned, b0.transactions, b0.blockType
FROM blocks b0
"""


def get_chain(database: Database, trace_from: int) -> Generator[Block, None, None]:
    for values in database.execute(QUERY_TRACE_CHAIN, (trace_from,)):
        yield Block.from_values(values)


def get_blocks(database: Database) -> Generator[Block, None, None]:
    for values in database.execute(QUERY_GET_BLOCKS):
        yield Block.from_values(values)


def get_all_blocks(database: Database) -> Generator[Block, None, None]:
    for values in database.execute(QUERY_GET_BLOCKS):
        yield Block.from_values(values)
