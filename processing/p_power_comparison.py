import sys

from model import load_blocks, load_label, load_nodes
from database import Database
from plot import plot_power_per_block_comparison


if __name__ == "__main__":
    if len(sys.argv) < 2:
        sys.exit(1)

    blocks_for_comparison = {}
    nodes_for_comparison = {}
    for i in range(1, len(sys.argv)):
        with Database(sys.argv[i]) as database:
            label = load_label(database)
            blocks = load_blocks(database)
            nodes = load_nodes(database)
            blocks_for_comparison[label] = blocks
            nodes_for_comparison[label] = nodes

    plot_power_per_block_comparison(blocks_for_comparison, nodes_for_comparison)
