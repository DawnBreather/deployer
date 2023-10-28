package main

import (
	"fmt"
	"time"
)

func main() {
	var i = 0
	for {
		i++
		fmt.Println(i)
		time.Sleep(3 * time.Second)
	}
}
