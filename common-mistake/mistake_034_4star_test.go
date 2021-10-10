package common_mistake

import (
	"fmt"
	"testing"
	"time"
)

// 使用了值为 nil 的 channel
// (1) 在一个值为 nil 的 channel 上发送和接收数据将永久阻塞
func TestMistake_034(t *testing.T) {
	wrong034()
	right034()
}

// 会导致死锁
func wrong034() {
	var ch chan int // 未初始化，值为 nil
	for i := 0; i < 3; i++ {
		go func(i int) {
			ch <- i
		}(i)
	}
	fmt.Println("Result: ", <-ch)
	time.Sleep(2 * time.Second)
}

// 利用这个死锁的特性，可以用在 select 中动态的打开和关闭 case 语句块：
func right034() {
	inCh := make(chan int)
	outCh := make(chan int)

	go func() {
		var in <-chan int = inCh
		var out chan<- int
		var val int

		for {
			select {
			case out <- val:
				println("--------")
				out = nil
				in = inCh
			case val = <-in:
				println("++++++++++")
				out = outCh
				in = nil
			}
		}
	}()

	go func() {
		for r := range outCh {
			fmt.Println("Result: ", r)
		}
	}()

	time.Sleep(0)
	inCh <- 1
	inCh <- 2
	time.Sleep(3 * time.Second)
}
