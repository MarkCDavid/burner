import sys

from model import load_blocks, load_label
from database import Database
from plot import plot_compare_production_times_typeless


if __name__ == "__main__":
    if len(sys.argv) < 2:
        sys.exit(1)

    blocks_for_comparison = {}
    for i in range(1, len(sys.argv)):
        with Database(sys.argv[i]) as database:
            label = load_label(database)
            blocks = load_blocks(database)
            blocks_for_comparison[label] = blocks

    plot_compare_production_times_typeless(blocks_for_comparison)
