from model import Block, BlockMiningAverages, PPOBPricing
from typing import Generator
from database import Database


QUERY_PPOB_CONSENSUS = """
SELECT 
    *
FROM pricing_proof_of_burn_burn_consensus
ORDER BY timestamp DESC;
"""


def query_ppob_pricing(
    database: Database,
) -> Generator[PPOBPricing, None, None]:
    for values in database.execute(QUERY_PPOB_CONSENSUS):
        yield PPOBPricing.from_values(values)
