package chapter16_Semaphore

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"
	"testing"
	"time"

	"golang.org/x/sync/semaphore"
)

// TestSemaphore_1
// (1)信号量
//
//		信号量分为 计数信号量 和 二进位信号量；有时候互斥锁也会使用二进位信号量来实现。
//		应用场景：一般用信号量保护一组资源，比如数据库连接池、一组客户端的连接、几个打印机资源，等等。
//	 sync.Mutex 中使用信号量来控制goroutine的阻塞和唤醒；不过并未对外暴露。
//			func runtime_Semacquire(s *uint32)
//			func runtime_SemacquireMutex(s *uint32, lifo bool, skipframes int)
//			func runtime_Semrelease(s *uint32, handoff bool, skipframes int)
//
// (2)扩展包中的信号量Weighted
//
//	扩展库：https://pkg.go.dev/golang.org/x
//	扩展库中的信号量Weighted：https://pkg.go.dev/golang.org/x/sync/semaphore
//	实现原理：互斥锁+List
//		type Weighted struct {
//			size    int64         // 最大资源数
//			cur     int64         // 当前已被使用的资源
//			mu      sync.Mutex    // 互斥锁，对字段的保护
//			waiters list.List     // 等待队列
//		}
//
// (3)常见错误
//
//	a.请求了资源，但是忘记释放它。
//	b.释放了从未请求的资源。
//	c.长时间持有一个资源，即使不需要它。
//	d.不持有一个资源，却直接使用它。
//	e.请求的资源数超过最大资源数。
func TestSemaphore_1(t *testing.T) {
	var (
		maxWorkers = runtime.GOMAXPROCS(0)                    // worker数量
		sema       = semaphore.NewWeighted(int64(maxWorkers)) // 信号量
		task       = make([]int, maxWorkers*4)                // 任务数，是worker的四倍
	)

	ctx := context.Background()

	for i := range task {
		// 如果没有worker可用，会阻塞在这里，直到某个worker被释放
		if err := sema.Acquire(ctx, 1); err != nil {
			break
		}

		// 启动worker goroutine
		go func(i int) {
			defer sema.Release(1)
			time.Sleep(100 * time.Millisecond) // 模拟一个耗时操作
			task[i] = i + 1
		}(i)
	}

	// 一个特别的技巧：请求所有的worker,这样能确保前面的worker都执行完
	if err := sema.Acquire(ctx, int64(maxWorkers)); err != nil {
		log.Printf("获取所有的worker失败: %v", err)
	}

	fmt.Println(task)
}

// TestSemaphore_2
// (4)其他信号量的实现
//
//	a.实现方式一：channel（例如：使用一个buffer为n的Channel很容易实现信号量）
//		Q:使用channel实现信号量非常简单，而且也能应对大部分的信号量的场景，为什么官方扩展库的信号量的实现不采用这种方法呢？
//		A:原因未知。但是官方的实现方式有这样一个功能：它可以一次请求多个资源，这是通过Channel实现的信号量所不具备的。
//	b.实现方式之二：https://github.com/marusama/semaphore
//		特点：一个可以动态更改资源容量的信号量。如果你的资源数量并不是固定的，而是动态变化的，那么建议考虑这个信号量库。
//
// (5)Q:你能用Channel实现信号量并发原语吗？你能想到几种实现方式？
//
//	A:至少两种，写入channel算获取，自己读取channel算获取。
//
// (6)Q:为什么信号量的资源数设计成int64而不是uint64呢？
//
//	A:防止错误获取或者释放信号量时，出现负数溢出到无穷大的问题。如果溢出到无穷大，就会让信号量失效，从而导致被保护资源更大规模的破坏。
func TestSemaphore_2(t *testing.T) {
	// a.实现方式一：channel（例如：使用一个buffer为n的Channel很容易实现信号量）
	_ = NewSemaphore(10)
}

// Semaphore 数据结构，并且还实现了Locker接口
type semaphore2 struct {
	sync.Locker
	ch chan struct{}
}

// 创建一个新的信号量
func NewSemaphore(capacity int) sync.Locker {
	if capacity <= 0 {
		capacity = 1 // 容量为1就变成了一个互斥锁
	}
	return &semaphore2{ch: make(chan struct{}, capacity)}
}

// 请求一个资源
func (s *semaphore2) Lock() {
	s.ch <- struct{}{}
}

// 释放资源
func (s *semaphore2) Unlock() {
	<-s.ch
}
