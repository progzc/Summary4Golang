package go_core_36

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

// TestExecImp 让主协程等待其他协程
func TestExecImp(t *testing.T) {
	num := 10
	sign := make(chan struct{}, num)
	for i := 0; i < num; i++ {
		// 引入局部变量,保证不会输出很多10
		a := i
		go func() {
			fmt.Println(a)
			sign <- struct{}{}
		}()
	}

	// 办法2。
	for j := 0; j < num; j++ {
		<-sign
	}
}

// TestExecImp2 保证顺序输出
func TestExecImp2(t *testing.T) {
	var count uint32
	trigger := func(i uint32, fn func()) {
		// 自旋
		for {
			if n := atomic.LoadUint32(&count); n == i {
				fn()
				atomic.AddUint32(&count, 1)
				break
			}
			time.Sleep(time.Nanosecond)
		}
	}
	for i := uint32(0); i < 10; i++ {
		// 增加入参,保证不会输出很多10
		go func(i uint32) {
			fn := func() {
				fmt.Println(i)
			}
			trigger(i, fn)
		}(i)
	}
	trigger(10, func() {})
}
