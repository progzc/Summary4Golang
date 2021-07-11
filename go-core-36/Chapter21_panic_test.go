package go_core_36

import (
	"errors"
	"fmt"
	"testing"
)

func caller() {
	fmt.Println("Enter function caller.")
	panic(errors.New("something wrong")) // 正例。
	panic(fmt.Println)                   // 反例。
	fmt.Println("Exit function caller.")
}

func TestPanic(t *testing.T) {
	fmt.Println("Enter function main.")
	caller()
	fmt.Println("Exit function main.")
}
