package common_mistake

import (
	"fmt"
	"testing"
)

// 从 panic 中恢复
// (1) 在一个 defer 延迟执行的函数中调用 recover() ，它便能捕捉/中断 panic
// (2) recover() 仅在 defer 执行的函数中调用才会生效。
func TestMistake_040(t *testing.T) {
	wrong040_1()
	wrong040_2()
	right040()
}

func wrong040_1() {
	recover()         // 什么都不会捕捉
	panic("not good") // 发生 panic，主程序退出
	recover()         // 不会被执行
	println("ok")
}

func wrong040_2() {
	defer func() {
		// (2) recover() 仅在 defer 执行的函数中调用才会生效。
		doRecover()
	}()
	panic("not good")
}

func doRecover() {
	fmt.Println("recobered: ", recover())
}

func right040() {
	defer func() {
		fmt.Println("recovered: ", recover())
	}()
	panic("not good")
}
