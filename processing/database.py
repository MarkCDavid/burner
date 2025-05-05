import os
import sqlite3


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

    def execute(self, query, params=None):
        if self._connection is None:
            raise ConnectionError()

        cursor = self._connection.cursor()
        cursor.execute(query, params or ())
        for row in cursor:
            yield row
        cursor.close()
