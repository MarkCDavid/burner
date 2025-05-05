from model import LeafBlock
from typing import Generator
from database import Database

QUERY_GET_LEAF_BLOCKS = """
WITH RECURSIVE chain AS (
  SELECT b0.id, b0.previousBlockId, b0.depth, b0.abandoned
  FROM blocks b0
  WHERE b0.previousBlockId = 0 AND b0.abandoned = 0

  UNION ALL

  SELECT b1.id, b1.previousBlockId, b1.depth, b1.abandoned
  FROM blocks b1
  JOIN chain c0 ON b1.previousBlockId = c0.id
  WHERE b1.abandoned = 0
)
SELECT c1.id, c1.depth
FROM chain c1
LEFT JOIN blocks b2 ON c1.id = b2.previousBlockId AND b2.abandoned = 0
WHERE b2.id IS NULL;
"""


def get_leaf_blocks(database: Database) -> Generator[LeafBlock, None, None]:
    for values in database.execute(QUERY_GET_LEAF_BLOCKS):
        yield LeafBlock.from_values(values)


def get_deepest_leaf(database: Database) -> LeafBlock:
    result = None
    for leaf in get_leaf_blocks(database):
        if result is None:
            result = leaf
            continue

        if result.depth < leaf.depth:
            result = leaf

    if result is None:
        raise ValueError()

    return result
