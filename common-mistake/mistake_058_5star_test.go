package common_mistake

import (
	"runtime"
	"testing"
)

// 优先调度
// (1) 你的程序可能出现一个 goroutine 在运行时阻止了其他 goroutine 的运行，比如程序中有一个不让调度器运行的 for 循环
// (2) 调度器会在 GC、Go 声明、阻塞 channel、阻塞系统调用和锁操作后再执行，也会在非内联函数调用时执行
//     添加 -m 参数来分析 for 代码块中调用的内联函数: go run -gcflags -m XXX.go
// (3) 可以使用runtime 包中的 Gosched() 来 手动启动调度器
func TestMistake_058(t *testing.T) {
	wrong058()
	//right058_1()
	//right058_2()
	right058_3()
}

func wrong058() {
}

func right058_1() {
	done := false
	go func() {
		done = true
	}()
	for !done {
	}
	println("done !")
}

func right058_2() {
	done := false
	go func() {
		done = true
	}()
	for !done {
		println("not done !") // 并不内联执行
	}
	println("done !")
}

func right058_3() {
	done := false
	go func() {
		done = true
	}()
	for !done {
		runtime.Gosched()
	}
	println("done !")
}
