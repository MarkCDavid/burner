from typing import Dict
import pandas as pd
import matplotlib.pyplot as plt
import numpy as np
from utility import trace_longest_chain, filter_abandoned
from tqdm import tqdm

tqdm.pandas()


def rebase_on(df: pd.DataFrame, field: str, on: pd.Timestamp):
    df[field] = df[field].apply(lambda t: on + pd.to_timedelta(t, unit="s"))
    return df


def calculate_mining_time(df: pd.DataFrame):
    df["miningTime"] = (df["finishedAt"] - df["startedAt"]).dt.total_seconds()
    df = df.sort_values("startedAt")
    return df


def calculate_production_time(df: pd.DataFrame, label: str) -> pd.DataFrame:
    df = df.sort_values("startedAt").copy()

    def compute_time(row):
        if row["blockType"] == 0 or "SlimCoin" in label:
            return (row["finishedAt"] - row["startedAt"]).total_seconds()
        elif row["blockType"] == 1 and pd.notnull(row["previousFinishedAt"]):
            return (row["finishedAt"] - row["previousFinishedAt"]).total_seconds()
        return None

    df["productionTime"] = df.apply(compute_time, axis=1)
    return df


def count_forks(blocks: pd.DataFrame) -> int:
    one_month_seconds = 30 * 24 * 60 * 60

    blocks = blocks[
        (blocks["abandoned"] == 0) & (blocks["startedAt"] >= one_month_seconds)
    ]

    forks_per_depth = blocks.groupby("depth")["previousBlockId"].nunique()
    total_forks = (forks_per_depth > 1).sum()

    return total_forks


def roll_window(df: pd.DataFrame, field: str, window: int = 1000):
    return df[field].rolling(window=window).mean()


def roll_window_timed(df: pd.DataFrame, field: str, window: str = "10min"):
    return df.set_index("startedAt")[field].rolling(window).mean()


def plot_monthly_lines(start_time: pd.Timestamp, end_time: pd.Timestamp):
    ticks = pd.date_range(start=start_time, end=end_time, freq="MS")

    for tick in ticks:
        plt.axvline(
            x=tick,
            color="black",
            linestyle="--",
            linewidth=1,
            alpha=0.5,
        )


def plot_target_line(target: int, label: str):
    plt.axhline(
        y=target,
        color="red",
        linestyle="--",
        linewidth=1,
        alpha=0.8,
        label=label,
    )


def plot_production_times(
    blocks: pd.DataFrame,
    label: str,
    window: int = 2016,
):
    start_time = pd.Timestamp("2024-01-01")
    end_time = pd.Timestamp("2025-01-01")

    blocks_longest = trace_longest_chain(blocks)
    blocks_longest = rebase_on(blocks_longest, "startedAt", start_time)
    blocks_longest = rebase_on(blocks_longest, "finishedAt", start_time)
    blocks_longest = rebase_on(blocks_longest, "previousFinishedAt", start_time)
    blocks_longest = blocks_longest[blocks_longest["startedAt"] < end_time]
    blocks_longest = calculate_production_time(blocks_longest, label)
    blocks_0 = blocks_longest[blocks_longest["blockType"] == 0]
    blocks_1 = blocks_longest[blocks_longest["blockType"] == 1]

    rolling_0 = roll_window(blocks_0, "productionTime", window)
    rolling_1 = roll_window(blocks_1, "productionTime", window // 16)

    plt.figure(figsize=(12, 6))

    plt.title(f"Rolling Avg. Block Production Time by Type ({label})")
    plt.xlabel("Simulation Time")
    plt.ylabel("Rolling Avg. Production Time (s)")

    plot_monthly_lines(start_time, end_time)
    plot_target_line(600, "Target: 600s")

    plt.plot(blocks_0["startedAt"], rolling_0, label="Proof of Work")
    plt.plot(blocks_1["startedAt"], rolling_1, label="Proof of Burn")

    plt.ylim(bottom=0)
    plt.legend()
    plt.tight_layout()
    plt.show()


def plot_production_times_typeless(
    blocks: pd.DataFrame,
    label: str,
    window: int = 2016,
):
    start_time = pd.Timestamp("2024-01-01")
    end_time = pd.Timestamp("2025-01-01")

    blocks_all = filter_abandoned(blocks)
    blocks_all = rebase_on(blocks_all, "startedAt", start_time)
    blocks_all = rebase_on(blocks_all, "finishedAt", start_time)
    blocks_all = rebase_on(blocks_all, "previousFinishedAt", start_time)
    blocks_all = calculate_production_time(blocks_all, label)
    blocks_all_rolling = roll_window(blocks_all, "productionTime", window)

    blocks_longest = trace_longest_chain(blocks)
    blocks_longest = rebase_on(blocks_longest, "startedAt", start_time)
    blocks_longest = rebase_on(blocks_longest, "finishedAt", start_time)
    blocks_longest = rebase_on(blocks_longest, "previousFinishedAt", start_time)
    blocks_longest = calculate_production_time(blocks_longest, label)
    blocks_longest_rolling = roll_window(blocks_longest, "productionTime", window)

    plt.figure(figsize=(12, 6))

    plt.title(f"Rolling Avg. Block Production Time ({label})")

    plt.xlabel("Simulation Time")
    plt.ylabel("Rolling Avg. Production Time (s)")

    plot_monthly_lines(start_time, end_time)
    plot_target_line(600, "Target: 600s")

    plt.plot(blocks_all["startedAt"], blocks_all_rolling, label="All Blocks")
    plt.plot(blocks_longest["startedAt"], blocks_longest_rolling, label="Longest Chain")

    plt.ylim(bottom=0)

    plt.legend()
    plt.tight_layout()
    plt.show()


def plot_compare_production_times_typeless(
    datasets: Dict[str, pd.DataFrame],
    window: int = 2016,
):
    start_time = pd.Timestamp("2024-01-01")
    end_time = pd.Timestamp("2025-01-01")

    plt.figure(figsize=(12, 6))
    plt.title("Rolling Avg. Block Production Time (Comparison)")
    plt.xlabel("Simulation Time")
    plt.ylabel("Rolling Avg. Production Time (s)")

    plot_monthly_lines(start_time, end_time)
    plot_target_line(600, "Target: 600s")

    for label, blocks in datasets.items():
        blocks_longest = trace_longest_chain(blocks)
        blocks_longest = rebase_on(blocks_longest, "startedAt", start_time)
        blocks_longest = rebase_on(blocks_longest, "finishedAt", start_time)
        blocks_longest = rebase_on(blocks_longest, "previousFinishedAt", start_time)
        blocks_longest = calculate_mining_time(blocks_longest)
        blocks_rolling = roll_window(blocks_longest, "miningTime", window)

        plt.plot(blocks_longest["startedAt"], blocks_rolling, label=label)

    plt.ylim(bottom=0)
    plt.legend()
    plt.tight_layout()
    plt.show()


def plot_transactions_per_block(
    datasets: Dict[str, pd.DataFrame],
    window: int = 2016,
):
    start_time = pd.Timestamp("2024-01-01")
    end_time = pd.Timestamp("2025-01-01")

    plt.figure(figsize=(12, 6))
    plt.title(f"Rolling Avg. Transactions per Block")

    plt.xlabel("Simulation Time")
    plt.ylabel("Rolling Avg. Transactions")

    plot_monthly_lines(start_time, end_time)
    plot_target_line(3000, "Max TX per Block")

    for label, blocks in datasets.items():
        blocks_longest = trace_longest_chain(blocks)
        blocks_longest = rebase_on(blocks_longest, "startedAt", start_time)
        blocks_longest = rebase_on(blocks_longest, "finishedAt", start_time)
        blocks_longest = rebase_on(blocks_longest, "previousFinishedAt", start_time)
        blocks_longest_rolling = roll_window(blocks_longest, "transactions", window)

        plt.plot(
            blocks_longest["startedAt"],
            blocks_longest_rolling,
            label=label,
        )

    plt.ylim(bottom=0)
    plt.legend()
    plt.tight_layout()
    plt.show()


def plot_power_per_block(
    blocks: Dict[str, pd.DataFrame],
    nodes: Dict[str, pd.DataFrame],
    window: int = 2016,
):
    start_time = pd.Timestamp("2024-01-01")
    end_time = pd.Timestamp("2025-01-01")

    plt.figure(figsize=(12, 6))
    plt.title(f"Rolling Avg. Power per Block")

    plt.xlabel("Simulation Time")
    plt.ylabel("Rolling Avg. Power")

    plot_monthly_lines(start_time, end_time)

    for label, blocks in blocks.items():
        blocks_full = blocks.copy()
        blocks_full = blocks_full.set_index("id")
        blocks_full = rebase_on(blocks_full, "startedAt", start_time)
        blocks_full = rebase_on(blocks_full, "finishedAt", start_time)
        blocks_full = rebase_on(blocks_full, "previousFinishedAt", start_time)
        blocks_full = calculate_mining_time(blocks_full)

        nodes_longest = nodes[label]
        merged_longest = blocks_full.merge(
            nodes_longest, left_on="minedBy", right_on="id", how="left"
        )

        merged_longest["powerUsage"] = merged_longest.apply(
            lambda row: row["miningTime"]
            * (row["powerFull"] if row["blockType"] == 0 else row["powerIdle"]),
            axis=1,
        )

        merged_longest_pow = merged_longest[merged_longest["blockType"] == 0].copy()
        merged_longest_pob = merged_longest[merged_longest["blockType"] == 1].copy()

        merged_longest_rolling = roll_window(merged_longest, "powerUsage", window)
        merged_longest_pow_rolling = roll_window(
            merged_longest_pow, "powerUsage", window
        )
        merged_longest_pob_rolling = roll_window(
            merged_longest_pob, "powerUsage", window
        )

        plt.plot(
            merged_longest["startedAt"],
            merged_longest_rolling,
            label=label,
        )

        plt.plot(
            merged_longest_pow["startedAt"],
            merged_longest_pow_rolling,
            label=f"{label} - PoW",
        )

        plt.plot(
            merged_longest_pob["startedAt"],
            merged_longest_pob_rolling,
            label=f"{label} - PoB",
        )

    plt.ylim(bottom=0)
    plt.legend()
    plt.tight_layout()
    plt.show()


def plot_power_per_block_comparison(
    blocks: Dict[str, pd.DataFrame],
    nodes: Dict[str, pd.DataFrame],
):
    start_time = pd.Timestamp("2024-01-01")
    end_time = pd.Timestamp("2025-01-01")

    plt.figure(figsize=(12, 6))
    plt.title(f"Cummulative Power Usage")

    plt.xlabel("Simulation Time")
    plt.ylabel("Power")

    plot_monthly_lines(start_time, end_time)

    for label, blocks in blocks.items():
        # blocks_full = blocks.copy()
        blocks_full = blocks[blocks["abandoned"] == 0].copy()
        blocks_full = blocks_full.set_index("id")
        # blocks_full = trace_longest_chain(blocks)
        blocks_full = rebase_on(blocks_full, "startedAt", start_time)
        blocks_full = rebase_on(blocks_full, "finishedAt", start_time)
        blocks_full = rebase_on(blocks_full, "previousFinishedAt", start_time)
        blocks_full = calculate_production_time(blocks_full, label)

        nodes_longest = nodes[label]

        nodes_longest_indexed = nodes_longest.set_index("id")
        blocks_full["powerFull"] = blocks_full["minedBy"].map(
            nodes_longest_indexed["powerFull"]
        )
        blocks_full["powerIdle"] = blocks_full["minedBy"].map(
            nodes_longest_indexed["powerIdle"]
        )

        blocks_full["powerUsage"] = blocks_full.progress_apply(
            lambda row: row["productionTime"]
            * (row["powerFull"] if row["blockType"] == 0 else row["powerIdle"]),
            axis=1,
        )

        blocks_full = blocks_full.sort_values("startedAt")
        blocks_full["cumulativePower"] = blocks_full["powerUsage"].cumsum()

        plt.plot(
            blocks_full["startedAt"],
            blocks_full["cumulativePower"],
            label=label,
        )

    # plt.ylim(bottom=0)
    plt.ylim(0, 1.75e12)
    plt.legend()
    plt.tight_layout()
    plt.show()


def plot_power_per_transaction_comparison(
    blocks: Dict[str, pd.DataFrame],
    nodes: Dict[str, pd.DataFrame],
):
    start_time = pd.Timestamp("2024-01-01")
    end_time = pd.Timestamp("2025-01-01")

    plt.figure(figsize=(12, 6))
    plt.title(f"Cumulative Power Usage")

    plt.xlabel("Simulation Time")
    plt.ylabel("Power")

    plot_monthly_lines(start_time, end_time)

    for label, blocksDf in blocks.items():
        blocks_full = blocksDf.copy()
        blocks_full = blocks_full.set_index("id")
        blocks_full = rebase_on(blocks_full, "startedAt", start_time)
        blocks_full = rebase_on(blocks_full, "finishedAt", start_time)
        blocks_full = rebase_on(blocks_full, "previousFinishedAt", start_time)
        blocks_full = calculate_production_time(blocks_full, label)

        nodes_longest = nodes[label]

        nodes_longest_indexed = nodes_longest.set_index("id")
        blocks_full["powerFull"] = blocks_full["minedBy"].map(
            nodes_longest_indexed["powerFull"]
        )
        blocks_full["powerIdle"] = blocks_full["minedBy"].map(
            nodes_longest_indexed["powerIdle"]
        )

        blocks_full["powerUsage"] = blocks_full.progress_apply(
            lambda row: row["productionTime"]
            * (row["powerFull"] if row["blockType"] == 0 else row["powerIdle"]),
            axis=1,
        )
        blocks_full["tx"] = blocks_full.progress_apply(
            lambda row: row["transactions"] if row["abandoned"] == 0 else 0,
            axis=1,
        )
        # blocks_full["pptx"] = blocks_full["powerUsage"] / blocks_full["tx"]

        blocks_full["pptx"] = blocks_full.apply(
            lambda row: row["powerUsage"] / row["tx"] if row["tx"] > 0 else np.nan,
            axis=1,
        )

        blocks_full = blocks_full.sort_values("startedAt")
        blocks_full_rolled = roll_window_timed(blocks_full, "pptx", "20160min")

        plt.plot(
            blocks_full["startedAt"],
            blocks_full_rolled,
            label=f"{label} - Rolling Avg Power/Tx",
        )

    plt.ylim(bottom=0)
    plt.legend()
    plt.tight_layout()
    plt.show()


def plot_rolling_avg_over_time(
    blocks: dict[str, pd.DataFrame],
    window: int = "20160min",
):
    start_time = pd.Timestamp("2024-01-01")
    end_time = pd.Timestamp("2025-01-01")

    plt.figure(figsize=(12, 6))
    plt.title(f"Cummulative Power Usage")

    plt.xlabel("Simulation Time")
    plt.ylabel("Power")

    plot_monthly_lines(start_time, end_time)

    for label, blocks in blocks.items():
        blocks_full = blocks[blocks["abandoned"] == 0].copy()
        blocks_full = rebase_on(blocks_full, "startedAt", start_time)
        blocks_full = rebase_on(blocks_full, "finishedAt", start_time)
        blocks_full = rebase_on(blocks_full, "previousFinishedAt", start_time)
        blocks_full.sort_values("startedAt", inplace=True)
        blocks_full.set_index("startedAt", inplace=True)

        series = pd.Series(1, index=blocks_full.index)
        resampled = series.resample("10min").sum()
        rolling_count = resampled.rolling(window=window, min_periods=1).sum()

        plt.plot(rolling_count.index, rolling_count.values, label=label)

    plt.legend()
    plt.tight_layout()
    plt.show()
