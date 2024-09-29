package common_mistake

import (
	"fmt"
	"testing"
)

// 失败的类型断言
// (1) 在类型断言语句中，断言失败则会返回目标类型的“零值”，断言变量与原来变量混用可能出现异常情况：
func TestMistake_050(t *testing.T) {
	wrong050()
	right050()
}

func wrong050() {
	var data interface{} = "great"
	// 错误用法
	if data, ok := data.(int); ok {
		fmt.Println("[is an int], data: ", data)
	} else {
		fmt.Println("[not an int], data: ", data) // [isn't a int], data:  0
	}
}

func right050() {
	var data interface{} = "great"
	// 正确用法
	if res, ok := data.(int); ok {
		fmt.Println("[is an int], data: ", res)
	} else {
		fmt.Println("[not an int], data: ", data) // [not an int], data:  great
	}
}
