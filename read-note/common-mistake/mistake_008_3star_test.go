package common_mistake

import (
	"testing"
)

// 显式类型的变量无法使用 nil 来初始化
// nil 是 interface、function、pointer、map、slice 和 channel 类型变量的默认初始值。
// 但声明时不指定类型，编译器也无法推断出变量的具体类型。
func TestMistake_008(t *testing.T) {
	wrong008()
	right008()
}

func wrong008() {
	//var x = nil    // error: use of untyped nil
	//_ = x
}

func right008() {
	var x interface{} = nil
	_ = x
}
