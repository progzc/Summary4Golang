package main

import (
	"fmt"
	"time"
)

// main 使用go run -race来检测代码是否存在潜在的并发读写
func main() {
	a := 1
	go func() {
		a = 2
	}()
	a = 3
	fmt.Println("a is ", a)
	time.Sleep(2 * time.Second)
}
