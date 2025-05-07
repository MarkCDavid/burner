from model import Block, BlockMiningAverages, POWPricing, PPOBPricing
from typing import Generator
from database import Database


QUERY_POW_CONSENSUS = """
SELECT 
    *
FROM proof_of_work_consensus
ORDER BY timestamp ASC;
"""


def query_pow_pricing(
    database: Database,
) -> Generator[POWPricing, None, None]:
    for values in database.execute(QUERY_POW_CONSENSUS):
        yield POWPricing.from_values(values)
