package common_mistake

import (
	"testing"
)

// 在多行 array、slice、map 语句中缺少 , 号
// (1) 声明语句中 } 折叠到单行后，尾部的 , 不是必需的。
func TestMistake_021(t *testing.T) {
	wrong021()
	right021()
}

func wrong021() {
	//x := []int {
	//	1,
	//	2    // syntax error: unexpected newline, expecting comma or }
	//}
}

func right021() {
	y := []int{1, 2}
	z := []int{1, 2}
	_ = y
	_ = z
}
