package main

import (
	"database/sql"
	"fmt"

	"github.com/kaiya/gorm/tidb"
)

func main() {
	db, err := sql.Open("mysql", "")
	if err != nil {
		panic(err)
	}
	td := tidb.NewTiDB(db)
	err = td.QueryAll()
	if err != nil {
		fmt.Println("query error: ", err)
	}
	fmt.Println("done")
}
