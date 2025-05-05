from typing import List, Dict
import matplotlib.pyplot as plt
import pandas as pd
from model import Block, PPOBPricing


def plot_block_times_by_type(blocks: Dict[int, Block]):
    powStartedAt = []
    powTimes = []
    pobStartedAt = []
    pobTimes = []

    for block in blocks.values():
        if block.blockType == 0:
            powStartedAt.append(block.startedAt)
            powTimes.append(block.duration())
        if block.blockType == 1:
            pobStartedAt.append(block.startedAt)
            pobTimes.append(block.duration())

    plt.figure(figsize=(12, 6))
    plt.scatter(powStartedAt, powTimes, color="green", label="Proof of Work")
    plt.scatter(pobStartedAt, pobTimes, color="red", label="Proof of Burn")

    plt.xlabel("StartedAt")
    plt.ylabel("Block Production Time")
    plt.title("Block Production Time by Block Type")
    plt.legend()
    plt.grid(True)
    plt.tight_layout()
    plt.show()


def plot_block_times_by_type_smoothed(
    blocks: Dict[int, Block], consensus: List[PPOBPricing], window_size=1000, name="N/A"
):
    # Process blocks
    records = []
    for index, block in blocks.items():
        if block.blockType == 0:
            records.append(
                {
                    "timestamp": block.startedAt,
                    "duration": block.duration(),
                    "blockType": block.blockType,
                }
            )
        elif block.blockType == 1:
            records.append(
                {
                    "timestamp": block.startedAt,
                    "duration": block.duration_since_previous_type(blocks),
                    "blockType": block.blockType,
                }
            )

    df_blocks = pd.DataFrame(records).sort_values(by="timestamp")
    df_blocks["rolling_avg"] = df_blocks.groupby("blockType")["duration"].transform(
        lambda x: x.rolling(window=window_size, min_periods=window_size).mean()
    )

    # Process consensus pricing
    df_price = pd.DataFrame([vars(p) for p in consensus]).sort_values(by="timestamp")
    df_price = df_price.groupby("timestamp")["price"].mean().reset_index()
    df_price["rolling_avg"] = (
        df_price["price"].rolling(window=window_size, min_periods=window_size).mean()
    )

    # Plot
    fig, ax1 = plt.subplots(figsize=(12, 6))

    # Block durations
    for block_type, color, label in [
        (0, "green", "Proof of Work"),
        (1, "red", "Proof of Burn"),
    ]:
        subset = df_blocks[df_blocks["blockType"] == block_type]
        ax1.plot(subset["timestamp"], subset["rolling_avg"], color=color, label=label)

    ax1.set_xlabel("Timestamp")
    ax1.set_ylabel("Block Production Time (s)")
    ax1.axhline(y=600, color="blue", linestyle="--", linewidth=0.5, label="Target 600s")
    ax1.set_ylim(bottom=0)
    ax1.grid(True)
    ax1.legend(loc="upper left")

    # Price consensus (right y-axis)
    ax2 = ax1.twinx()
    ax2.plot(
        df_price["timestamp"],
        df_price["rolling_avg"],
        color="orange",
        label="PPOB Price (avg)",
    )
    ax2.set_ylabel("PPOB Consensus Price")
    ax2.legend(loc="upper right")

    plt.title("Block Production Time and PPOB Consensus Price (Smoothed)")
    manager = plt.gcf().canvas.manager
    if manager:
        manager.set_window_title(name)
    plt.tight_layout()
    plt.show()


#
# def plot_block_times_by_type_smoothed(
#     blocks: Dict[int, Block], consensus: List[PPOBPricing], window_size=1000, name="N/A"
# ):
#     records = []
#
#     for index, block in blocks.items():
#         if block.blockType == 0:
#             records.append(
#                 {
#                     "depth": block.depth,
#                     "duration": block.duration(),
#                     "blockType": block.blockType,
#                 }
#             )
#
#         if block.blockType == 1:
#             records.append(
#                 {
#                     "depth": block.depth,
#                     "duration": block.duration_since_previous_type(blocks),
#                     "blockType": block.blockType,
#                 }
#             )
#
#     df = pd.DataFrame(records).sort_values(by="depth")
#     df["rolling_avg"] = df.groupby("blockType")["duration"].transform(
#         lambda x: x.rolling(window=window_size, min_periods=window_size).mean()
#     )
#
#     plt.figure(figsize=(12, 6))
#
#     for block_type, color, label in [
#         (0, "green", "Proof of Work"),
#         (1, "red", "Proof of Burn"),
#     ]:
#         subset = df[df["blockType"] == block_type]
#         plt.plot(subset["depth"], subset["rolling_avg"], color=color, label=label)
#
#     max_depth = df["depth"].max()
#     for epoch in range(0, max_depth + 1, 2016):
#         plt.axvline(x=epoch, color="gray", linestyle="--", linewidth=0.5)
#
#     plt.axhline(y=600, color="blue", linestyle="--", linewidth=0.5)
#
#     plt.xlabel("Depth")
#     plt.ylabel(f"Rolling Avg Block Production Time (window={window_size})")
#     plt.title("Smoothed Block Production Time by Block Type")
#     plt.ylim(bottom=0)
#     plt.legend()
#     manager = plt.gcf().canvas.manager
#     if manager:
#         manager.set_window_title(name)
#     plt.tight_layout()
#     plt.show()
