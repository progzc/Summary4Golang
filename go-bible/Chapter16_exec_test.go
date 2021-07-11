package go_bible

import (
	"fmt"
	"testing"
	"time"
)

// TestExec 最初版本,bug很多
func TestExec(t *testing.T) {
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println(i)
		}()
	}
}

// TestExec2 保证协程执行完
func TestExec2(t *testing.T) {
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println(i)
		}()
	}
	time.Sleep(100)
}

// TestExec3 解决闭包问题
func TestExec3(t *testing.T) {
	for i := 0; i < 10; i++ {
		a := i
		go func() {
			fmt.Println(a)
		}()
	}
	time.Sleep(100)
}
