import pandas as pd
from tqdm import tqdm


def filter_abandoned(df: pd.DataFrame):
    df = df[df["abandoned"] == 0].copy()
    return df


def trace_longest_chain(blocks: pd.DataFrame):
    blocks_all = filter_abandoned(blocks)
    blocks_all = blocks_all.set_index("id")

    tip = blocks_all[blocks_all["depth"] == blocks_all["depth"].max()].sample(1)
    current_id = tip.index[0]
    target_depth = tip["depth"].iloc[0]

    chain = []
    previous_depth = target_depth
    with tqdm(total=target_depth + 1, desc="Tracing main chain", leave=False) as bar:
        while pd.notna(current_id):
            block = blocks_all.loc[current_id]
            chain.append(block)

            current_depth = block["depth"]
            bar.update(previous_depth - current_depth)
            previous_depth = current_depth

            current_id = (
                block["previousBlockId"]
                if block["previousBlockId"] in blocks_all.index
                else None
            )

    return pd.DataFrame(reversed(chain)).reset_index()
