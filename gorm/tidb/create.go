package tidb

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type TiDB struct {
	db *sql.DB
}

func NewTiDB(db *sql.DB) *TiDB {
	return &TiDB{
		db: db,
	}
}

/*
func (td *TiDB) InsertOne() error {
}
*/

func (td *TiDB) QueryAll() error {
	rows, err := td.db.Query("select User from user")
	if err != nil {
		return errors.Wrap(err, "query user")
	}
	for {
		var user string
		if rows.Next() {
			err := rows.Scan(&user)
			if err != nil {
				return errors.Wrap(err, "scan next")
			}
			fmt.Printf("Got one, User:%s\n", user)
		} else {
			break
		}
	}
	return nil
}
