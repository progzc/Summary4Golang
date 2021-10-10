package common_mistake

import (
	"fmt"
	"testing"
)

// 跳出 for-switch 和 for-select 代码块
// (1) 没有指定标签的 break 只会跳出 switch/select 语句，若不能使用 return 语句跳出的话，可为 break 跳出标签指定的代码块
func TestMistake_046(t *testing.T) {
	wrong046()
	right046()
}

func wrong046() {
}

func right046() {
loop:
	for {
		switch {
		case true:
			fmt.Println("breaking out...")
			//break // 死循环，一直打印 breaking out...
			break loop
		}
	}
	fmt.Println("out...")
}
