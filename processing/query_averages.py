from model import Block, BlockMiningAverages
from typing import Generator
from database import Database


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


def query_block_production_times(
    database: Database,
) -> Generator[BlockMiningAverages, None, None]:
    for values in database.execute(QUERY_BLOCK_PRODUCTION_TIMES):
        yield BlockMiningAverages.from_values(values)
