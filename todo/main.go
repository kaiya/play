package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kaiya/play/todo/dao"
)

var (
	sqlHost = flag.String("sqlUser", "gateway01.ap-northeast-1.prod.aws.tidbcloud.com", "host for login mysql")
	sqlPort = flag.Int("sqlPort", 4000, "port for  mysql")
	sqlDb   = flag.String("sqlDb", "todo", "username for login mysql")
	sqlUser = flag.String("sqlUser", "2XnQNnZnVfDF5Qq.root", "username for login mysql")
	sqlPass = flag.String("sqlPass", "", "password for mysql")
)

func main() {
	flag.Parse()
	err := SetFlagsFromEnvironment()
	if err != nil {
		panic(err)
	}
	if *sqlPass == "" {
		*sqlPass = os.Getenv("sqlPass")
	}

	if *sqlUser == "" || *sqlPass == "" {
		panic("sqlPass is empty")
	}
	dao.Init(*sqlUser, *sqlPass, *sqlHost, *sqlDb, *sqlPort)
	router := gin.Default()
	api := router.Group("/api/v1/todos")
	{
		api.POST("/", dao.CreateTodo)
		api.POST("/bind", dao.BindCreate)
		api.GET("/", dao.FetchAllTodo)
		api.GET("/:id", dao.FetchSingleTodo)
		api.PUT("/:id", dao.UpdateTodo)
		api.DELETE("/:id", dao.DeleteTodo)
	}
	router.Run()
}

func SetFlagsFromEnvironment() (err error) {
	flag.VisitAll(func(f *flag.Flag) {
		name := strings.ToUpper(strings.Replace(f.Name, "-", "_", -1))
		if value, ok := os.LookupEnv(name); ok {
			err2 := flag.Set(f.Name, value)
			if err2 != nil {
				err = fmt.Errorf("failed setting flag from environment: %w", err2)
			}
		}
	})

	return
}
