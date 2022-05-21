package common_mistake

import (
	"testing"
)

// 不小心覆盖了变量
// 对从动态语言转过来的开发者来说，简短声明很好用，这可能会让人误会 := 是一个赋值操作符。
func TestMistake_007(t *testing.T) {
	right007()
}

func right007() {
	x := 1
	println(x) // 1
	{
		println(x) // 1
		// 新的 x 变量的作用域只在代码块内部
		// 这是 Go 开发者常犯的错，而且不易被发现
		// 解决办法:
		// (1) 使用 vet 工具来诊断这种变量覆盖，Go 默认不做覆盖检查，添加 -shadow 选项来启用
		//	   例如：go_knowledge tool vet -shadow 6_escape_1.go_knowledge
		// (2) vet 不会报告全部被覆盖的变量，可以使用 go_knowledge-nyet 来做进一步的检测
		//	   例如：$GOPATH/bin/go_knowledge-nyet 6_escape_1.go_knowledge
		x := 2
		println(x) // 2
	}
	println(x) // 1
}
