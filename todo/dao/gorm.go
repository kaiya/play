package dao

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init(sqlUser, sqlPass, sqlHost, sqlDb string, sqlPort int) {
	var err error
	db, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", sqlUser, sqlPass, sqlHost, sqlPort, sqlDb)))
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&todoModel{})
}
