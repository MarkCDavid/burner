import os
import sqlite3
import pandas as pd
from tqdm import tqdm


class Database:
    def __init__(self, path):
        self._path = path
        self._name = os.path.basename(path)
        self._connection = None

    def __enter__(self):
        self._connection = sqlite3.connect(self._path)
        return self

    def __exit__(self, exc_type, exc_value, traceback):
        if self._connection:
            self._connection.close()
            self._connection = None

    def get_table_count(self, table):
        if self._connection is None:
            raise ConnectionError()

        cursor = self._connection.cursor()
        cursor.execute(f"SELECT COUNT(*) FROM {table}")

        return int(cursor.fetchone()[0])

    def execute_df(
        self,
        table,
        columns,
    ):
        if self._connection is None:
            raise ConnectionError()

        count = self.get_table_count(table)

        print(f"database: fetching from {table}...")

        cursor = self._connection.cursor()
        cursor.execute(f"SELECT * FROM {table} ORDER BY id ASC")

        rows = list(
            tqdm(
                cursor,
                total=count,
                desc=f"fetching from {table}...",
                leave=False,
            )
        )
        df = pd.DataFrame(rows, columns=columns)
        return df
