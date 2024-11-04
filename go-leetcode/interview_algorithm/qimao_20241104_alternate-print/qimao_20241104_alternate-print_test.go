package qimao_20241104_alternate_print

import (
	"fmt"
	"testing"
)

// 七猫（一面）
// 两个协程交替打印输出1-10

// TestAlternatePrint
func TestAlternatePrint(t *testing.T) {
	chs := make([]chan int, 2)
	for i := 0; i < len(chs); i++ {
		chs[i] = make(chan int)
	}

	done := make(chan struct{})
	start, target := 1, 10
	go func() {
		for {
			select {
			case x := <-chs[0]:
				if x > target {
					close(done)
					return
				}
				fmt.Printf("goroutine1: %d\n", x)
				y := x + 1
				chs[1] <- y
			case <-done:
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case x := <-chs[1]:
				if x > target {
					close(done)
					return
				}
				fmt.Printf("goroutine2: %d\n", x)
				y := x + 1
				chs[0] <- y
			case <-done:
				return
			}
		}
	}()

	chs[0] <- start
	select {
	case <-done:
		return
	}
}
