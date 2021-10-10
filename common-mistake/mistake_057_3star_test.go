package common_mistake

import (
	"testing"
	"time"
)

// 读写操作的重新排序
// (1) Go 可能会重排一些操作的执行顺序，可以保证在一个 goroutine 中操作是顺序执行的，但不保证多 goroutine 的执行顺序
func TestMistake_057(t *testing.T) {
	wrong057()
	right057()
}

func wrong057() {
}

func right057() {
	var a, b int

	u1 := func() {
		a = 1
		b = 2
	}

	u2 := func() {
		a = 3
		b = 4
	}

	p := func() {
		println(a)
		println(b)
	}

	go u1() // 多个 goroutine 的执行顺序不定
	go u2()
	go p()
	time.Sleep(1 * time.Second)
}
