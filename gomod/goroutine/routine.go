package goroutine

import (
	"fmt"
	"time"
)

func Sleep(s string) {
	var i int = 0
	for {
		i += 1
		fmt.Println(i)
		fmt.Println(s)
		time.Sleep(time.Second)
	}
}
