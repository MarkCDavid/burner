import pandas as pd
from database import Database


def load_label(db: Database) -> str:
    columns = ["id", "label"]
    df = db.execute_df("label", columns)
    if df.empty:
        raise Exception("No label in database.")

    return df.sample(1)["label"].iloc[0]


def load_blocks(db: Database) -> pd.DataFrame:
    columns = [
        "id",
        "previousBlockId",
        "minedBy",
        "depth",
        "startedAt",
        "finishedAt",
        "previousFinishedAt",
        "abandoned",
        "transactions",
        "blockType",
    ]

    df = db.execute_df("blocks", columns)
    return df


def load_nodes(db: Database) -> pd.DataFrame:
    columns = [
        "id",
        "powerFull",
        "powerIdle",
        "transactions",
    ]
    return db.execute_df("nodes", columns)


def load_pricing_proof_of_burn_burn_consensus(db: Database) -> pd.DataFrame:
    columns = [
        "id",
        "timestamp",
        "nodeId",
        "currentlyAt",
        "price",
    ]
    return db.execute_df("pricing_proof_of_burn_burn_consensus", columns)


def load_pricing_proof_of_burn_burn_transaction(db: Database) -> pd.DataFrame:
    columns = [
        "id",
        "nodeId",
        "burnedAt",
        "burnedFor",
    ]
    return db.execute_df("pricing_proof_of_burn_burn_transaction", columns)


def load_proof_of_work_consensus(db: Database) -> pd.DataFrame:
    columns = [
        "id",
        "timestamp",
        "nodeId",
        "difficulty",
        "eventType",
    ]
    return db.execute_df("proof_of_work_consensus", columns)


def load_razer_proof_of_burn_consensus(db: Database) -> pd.DataFrame:
    columns = [
        "id",
        "timestamp",
        "nodeId",
        "chance",
        "eventType",
    ]
    return db.execute_df("razer_proof_of_burn_consensus", columns)


def load_slimcoin_proof_of_burn_consensus(db: Database) -> pd.DataFrame:
    columns = [
        "id",
        "timestamp",
        "nodeId",
        "chance",
        "eventType",
    ]
    return db.execute_df("slimcoin_proof_of_burn_consensus", columns)
