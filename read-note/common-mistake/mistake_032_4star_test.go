package common_mistake

import (
	"fmt"
	"testing"
	"time"
)

// 向无缓冲的 channel 发送数据，只要 receiver 准备好了就会立刻返回
// (1) 只有在数据被 receiver 处理时，sender 才不会阻塞。
// (2) 因运行环境而异，在 sender 发送完数据后，receiver 的 goroutine 可能没有足够的时间处理下一个数据。如：
func TestMistake_032(t *testing.T) {
	wrong032()
	right032()
}

func wrong032() {
	ch := make(chan string)

	go func() {
		for m := range ch {
			fmt.Println("Processed:", m)
			time.Sleep(1 * time.Second) // 模拟需要长时间运行的操作
		}
	}()

	ch <- "cmd.1"
	ch <- "cmd.2" // 不会被接收处理
}

func right032() {
}
