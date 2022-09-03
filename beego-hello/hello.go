package main

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql" // import your used driver
)

type User struct {
	Id   int
	Name string `orm:"size(100)"`
}

func init() {
	// set default database
	orm.RegisterDataBase("default", "mysql", "root:kaiyaxiong123@tcp(127.0.0.1:3306)/beego")
	//register model
	orm.RegisterModel(new(User))
	//create table
	orm.RunSyncdb("default", false, true)
}

type UserController struct {
	web.Controller
}

func (uc *UserController) Get() {
	o := orm.NewOrm()
	user := &User{Name: "Kaiya"}
	//insert
	id, err := o.Insert(user)
	fmt.Printf("ID: %d, ERR:%v\n", id, err)
	//update
	user.Name = "xiong"
	num, err := o.Update(user)
	fmt.Printf("NUM: %d, ERR: %v\n", num, err)
	//read one
	u := &User{Id: user.Id}
	err = o.Read(u)
	fmt.Printf("ERR: %v\n", err)
	//read to the u object
	//delete
	// num, err = o.Delete(u)
	// fmt.Printf("NUM:%d, ERR: %v\n", num, err)
	uc.Ctx.WriteString("hello: " + u.Name + ", and Id: " + fmt.Sprint(u.Id))
}
func mainha() {
	web.Router("/hello", &UserController{})
	web.Run()
}
