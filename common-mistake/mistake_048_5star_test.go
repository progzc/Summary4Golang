package common_mistake

import (
	"fmt"
	"testing"
)

// defer 函数的参数值
// (1) 对 defer 延迟执行的函数，它的参数会在声明时候就会求出具体值，而不是在执行时才求值
func TestMistake_048(t *testing.T) {
	wrong048()
	right048()
}

func wrong048() {
}

func right048() {
	var i = 1
	defer fmt.Println("result: ", func() int { return i * 2 }()) // result:  2
	i++
}
