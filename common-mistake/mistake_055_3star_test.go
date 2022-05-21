package common_mistake

import (
	"testing"
)

// 堆栈变量
// (1) 你并不总是清楚你的变量是分配到了堆还是栈。
//     在 C++ 中使用 new 创建的变量总是分配到堆内存上的，但在 Go 中即使使用 new()、make() 来创建变量，变量为内存分配位置依旧归 Go 编译器管。
// (2) Go 编译器会根据变量的大小及其 "escape analysis" 的结果来决定变量的存储位置，故能准确返回本地变量的地址，这在 C/C++ 中是不行的。
// (3) 在 go_knowledge build 或 go_knowledge run 时，加入 -m 参数，能准确分析程序的变量分配位置：
//     例：go_knowledge run -gcflags -m XXXX.go_knowledge
func TestMistake_055(t *testing.T) {
	wrong055()
	right055()
}

func wrong055() {
}

func right055() {
}
