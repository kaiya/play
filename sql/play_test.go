package playsql

import (
	"database/sql"
	"net"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"gitlab.momoso.com/mms2/utils/lg"
)

func TestSql(t *testing.T) {
	// sql raw data encoding
	/*
		conn, err := mysql.NewConnector(&mysql.Config{
			User:                 "2XnQNnZnVfDF5Qq.root",
			Passwd:               "kaiyaxiong123",
			ParseTime:            true,
			DBName:               "todo",
			Net:                  "tcp4",
			Addr:                 "gateway01.ap-northeast-1.prod.aws.tidbcloud.com:4000",
			AllowNativePasswords: true,
			Collation:            "utf8mb4_general_ci",
			Loc:                  time.UTC,
		})
		lg.PanicError(err)
		db := sql.OpenDB(conn)
	*/
	db, err := sql.Open("mysql", "root:@tcp(localhost:4000)/todo?parseTime=true")
	lg.PanicError(err)
	rows, err := db.Query("select * from todo_models")
	lg.PanicError(err)
	defer rows.Close()
	// cached colomns
	// rows2, err := db.Query("select * from todo_models")
	for {
		if rows.Next() {
			rows.ColumnTypes()
			columns := make([]sql.RawBytes, 6)
			err = rows.Scan(&columns[0], &columns[1], &columns[2], &columns[3], &columns[4], &columns[5])
			if err != nil {
				lg.Errorf("scan error:%s", err)
			}
		} else {
			break
		}
	}
}

func TestMysqlTcp(t *testing.T) {
	tcpConn, err := net.Dial("tcp4", "localhost:3306")
	lg.PanicError(err)
	_ = tcpConn

}
