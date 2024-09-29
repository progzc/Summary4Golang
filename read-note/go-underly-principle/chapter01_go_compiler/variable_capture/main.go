package main

import "fmt"

// go_knowledge tool compile -m=2 main.go_knowledge | grep capturing
func main() {
	// a采用引用传递，b采用值传递
	a := 1
	b := 2
	go func() {
		fmt.Println(a, b)
	}()
	a = 99
}
