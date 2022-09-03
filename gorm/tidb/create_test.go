package tidb

import (
	"database/sql"
	"testing"
)

func Benchmark_create(b *testing.B) {
	db, err := sql.Open("mysql", "")
	if err != nil {
		panic(err)
	}
	td := NewTiDB(db)
	b.Run("b-create", func(b *testing.B) {
		td.QueryAll()
	})
}

func Test_query(t *testing.T) {
	t.Run("query", func(t *testing.T) {
		db, err := sql.Open("mysql", "")
		if err != nil {
			panic(err)
		}
		td := NewTiDB(db)
		err = td.QueryAll()
		if err != nil {
			t.Error(err)
		}
	})
}
