package go_core_36

import (
	"errors"
	"fmt"
	"testing"
)

// TestRecovery 不正确的用法
func TestRecovery(t *testing.T) {
	fmt.Println("Enter function main.") // 引发panic。
	panic(errors.New("something wrong"))
	p := recover()
	fmt.Printf("panic: %s\n", p)
	fmt.Println("Exit function main.")
}

// TestRecovery2 正确的用法
func TestRecovery2(t *testing.T) {
	fmt.Println("Enter function main.")
	defer func() {
		fmt.Println("Enter defer function.")
		if p := recover(); p != nil {
			fmt.Printf("panic: %s\n", p)
		}
		fmt.Println("Exit defer function.")
	}() // 引发panic。
	panic(errors.New("something wrong"))
	fmt.Println("Exit function main.")
}

// TestDefer defer的使用
func TestDefer(t *testing.T) {
	defer fmt.Println("first defer")
	for i := 0; i < 3; i++ {
		defer fmt.Printf("defer in for [%d]\n", i)
	}
	defer fmt.Println("last defer")
}
