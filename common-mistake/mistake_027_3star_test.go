package common_mistake

import (
	"fmt"
	"testing"
)

// 自增和自减运算
// (1) Go 特立独行，去掉了++/--的前置操作
// (2) ++/--只能作为运算符而非表达式
func TestMistake_027(t *testing.T) {
	wrong027()
	right027()
}

func wrong027() {
	//data := []int{1, 2, 3}
	//i := 0
	//++i            // syntax error: unexpected ++, expecting }
	//fmt.Println(data[i++])    // syntax error: unexpected ++, expecting :
}

func right027() {
	data := []int{1, 2, 3}
	i := 0
	i++
	fmt.Println(data[i]) // 2
}
