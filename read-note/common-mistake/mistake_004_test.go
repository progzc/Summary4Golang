package common_mistake

import (
	"testing"
)

// 简短声明的变量只能在函数内部使用
func TestMistake_004(t *testing.T) {
	wrong004()
	right004()
}

// 错误示例
// myvar := 1
func wrong004() {
}

// 正确示例
// var  myvar = 1
func right004() {
}
