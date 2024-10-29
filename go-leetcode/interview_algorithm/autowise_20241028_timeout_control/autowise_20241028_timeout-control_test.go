package autowise_20241028_timeout_control

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

// 仙途智能（一面）
// 从多个数据源并行的获取数据。过程中需要实现对于获取数据的超时控制。并用一个channel收集所有数据。

func TestHandle(t *testing.T) {
	Handle()
}

func Handle() {
	rand.Seed(time.Now().UnixNano())
	ch := make(chan string)
	done := make(chan struct{})
	once := sync.Once{}
	targetNum := 3

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	for i := 0; i < targetNum; i++ {
		go func(ctx context.Context, i int) {
			// 1.异步执行逻辑
			data := make(chan string)
			go func(ctx context.Context) {
				longTimeTask(i)
				data <- fmt.Sprintf("result from %d", i)
			}(ctx)
			// 2.控制自身超时、上游超时、处理完数据提前结束
			select {
			case <-time.After(time.Second * 2):
				fmt.Printf("sub goroutine timeout, lead sub(%d) to exit\n", i)
				once.Do(func() { close(done) })
				return
			case <-ctx.Done():
				fmt.Printf("main goroutine timeout, lead sub(%d) to exit\n", i)
				once.Do(func() { close(done) })
				return
			case x := <-data:
				ch <- x
				return
			}
		}(ctx, i)
	}

	for i := 0; i < targetNum; i++ {
		select {
		case result, ok := <-ch:
			if !ok {
				return
			} else {
				fmt.Printf("get result: %v\n", result)
			}
		case <-ctx.Done():
			fmt.Printf("main goroutine timeout to exit\n")
			return
		case <-done:
			fmt.Printf("main goroutine timeout signal to exit\n")
			return
		}
	}
}

func longTimeTask(i int) {
	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
}
