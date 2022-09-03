package dao

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init(sqlUser, sqlPass string) {
	var err error
	db, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(gateway01.ap-northeast-1.prod.aws.tidbcloud.com:4000)/todo", sqlUser, sqlPass)))
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&todoModel{})
}
