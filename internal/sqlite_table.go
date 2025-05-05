package internal

import (
	"database/sql"
	"reflect"

	"github.com/sirupsen/logrus"
)

type SQLiteTable[T any] struct {
	_sqlite   *SQLite
	_batching []T
	_insert   *sql.Stmt
}

func NewSQLiteTable[T any](
	sqlite *SQLite,
	createQuery string,
	insertQuery string,
) *SQLiteTable[T] {
	_, err := sqlite._database.Exec(createQuery)
	if err != nil {
		logrus.Fatal(err)
	}

	insert, err := sqlite._database.Prepare(insertQuery)
	if err != nil {
		logrus.Fatal(err)
	}

	return &SQLiteTable[T]{
		_sqlite:   sqlite,
		_batching: []T{},
		_insert:   insert,
	}
}

func (s *SQLiteTable[T]) Save(entity T) {
	s._batching = append(s._batching, entity)
	if len(s._batching) >= batchSize {
		s.Flush()
	}
}

func (s *SQLiteTable[T]) Flush() {
	if len(s._batching) == 0 {
		return
	}

	transaction, err := s._sqlite._database.Begin()
	if err != nil {
		logrus.Fatal(err)
	}

	statement := transaction.Stmt(s._insert)

	for _, element := range s._batching {
		_, err := statement.Exec(Slice(element)...)
		if err != nil {
			transaction.Rollback()
			logrus.Fatal(err)
		}
	}

	if err := transaction.Commit(); err != nil {
		logrus.Fatal(err)
	}

	s._batching = s._batching[:0]

}

func Slice(s any) []any {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Struct {
		panic("expected a struct")
	}

	result := make([]any, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		result[i] = v.Field(i).Interface()
	}
	return result
}
