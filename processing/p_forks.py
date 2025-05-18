import sys

from model import load_blocks, load_label
from database import Database
from plot import count_forks


if __name__ == "__main__":
    if len(sys.argv) < 2:
        sys.exit(1)

    blocks_for_comparison = {}
    nodes_for_comparison = {}
    for i in range(1, len(sys.argv)):
        with Database(sys.argv[i]) as database:
            label = load_label(database)
            blocks = load_blocks(database)
            forks = count_forks(blocks)

            print(f"Forks ({label}): {forks}")
