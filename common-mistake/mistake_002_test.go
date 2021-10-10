package common_mistake

import "testing"

// 错误示例
//var gvar int // 全局变量，声明不使用也可以

// 未使用的变量
// 如果在函数体代码中有未使用的变量，则无法通过编译，不过全局变量声明但不使用是可以的。
// 即使变量声明后为变量赋值，依旧无法通过编译，需在某处使用它
func TestMistake_002(t *testing.T) {
	// 错误示例
	//var one int     // error: one declared and not used
	//two := 2    // error: two declared and not used
	//var three int    // error: three declared and not used
	//three = 3

	// 正确示例
	var one int
	_ = one

	two := 2
	println(two)

	var three int
	one = three

	var four int
	four = four
}
