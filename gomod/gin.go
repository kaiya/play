package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	go func(url string) {
		resp, err := http.DefaultClient.Get(url)
		if err != nil {
			log.Fatalln("get error!")
		}
		fmt.Println(resp)
	}("https://kaiyai.com")
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"isGet":   true,
			"message": "world",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
