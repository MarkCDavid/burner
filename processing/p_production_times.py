import sys

from model import load_blocks, load_label
from database import Database
from plot import plot_production_times


if __name__ == "__main__":
    if len(sys.argv) != 2:
        sys.exit(1)

    with Database(sys.argv[1]) as database:
        label = load_label(database)
        blocks = load_blocks(database)
        plot_production_times(blocks, label, window=4000)
