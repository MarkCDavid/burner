import sys

from database import Database
from plot_block_times import plot_block_times_by_type, plot_block_times_by_type_smoothed
from query_pow import query_pow_pricing
from query_averages import query_block_production_times
from query_blocks import get_chain
from query_leaf import get_deepest_leaf, get_leaf_blocks
from query_ppob import query_ppob_pricing

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python script.py <path_to_sqlite_file>")
        sys.exit(1)

    with Database(sys.argv[1]) as database:
        blocks = {k.id: k for k in get_chain(database, get_deepest_leaf(database).id)}
        consensus = list(query_ppob_pricing(database))
        consensusPow = list(query_pow_pricing(database))
        # plot_block_times_by_type(blocks)
        plot_block_times_by_type_smoothed(
            blocks, consensusPow, consensus, 400, database._name
        )
        # for time in query_block_production_times(database):
        #     print(time)
