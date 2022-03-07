package chapter31_sync_WaitGroup_Once

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestWaitGroupOnce_1
// (1)需求:声明一个通道，使它的容量与我们手动启用的goroutine的数量相同，之后再利用这个通道，让主goroutine等待其他 goroutine的运行结束
//	常规方式: 方式一使用channel、方式二使用sync.WaitGroup（推荐）
// (2)Q:sync.WaitGroup类型值中计数器的值可以小于0吗?
//	  A:不可以,否则会引发panic
//	容易出现panic的场景:
//	  a.wg.Add方法传入负数(虽然wg.Add方法可以传入负数,但是绝对不允许这么做,这会直接导致程序panic)
//	  b.对Add方法的首次调用,与对它的Wait方法的调用是同时发起 (即两个方法的调用在同时启动的不同的goroutine中)
//	  c.在复用sync.WaitGroup时,未保证其计数周期的完整性
// (3)使用sync.WaitGroup的注意事项:
//	  a.使用sync.WaitGroup时可以复用的,弹药保证其技术周期性
//	  b.不要把Add方法和Wait方法的代码，放在不同的goroutine中执行。即,要杜绝对同一个WaitGroup值的两种操作的并发执行。
//	  c.禁止在函数中对sync.WaitGroup进行直接值复制,应该采用指针复制
// (4)Q:sync.Once类型值的Do方法是怎么保证只执行参数函数一次的?sync.WaitGroup类型的实现原理呢?
//	  A:
//	 	a.sync.Once的实现原理: uint32(值只有从0变为1的这种状态)+原子操作+互斥锁，本质时双重检查(即单例模式的double check机制)。
//		b.sync.WaitGroup的实现原理: 原子操作
// (5)Q:在使用WaitGroup值实现一对多的 goroutine 协作流程时，怎样才能让分发子任务的 goroutine 获得各个子任务的具体执行结果?
//	  A:可以考虑使用锁 + 容器（数组、切片或字典等），也可以考虑使用通道。
//	    另外，你或许也可以用上golang.org/x/sync/errgroup代码包中的程序实体，相应的文档在这里。
func TestWaitGroupOnce_1(t *testing.T) {
	// (1)需求:声明一个通道，使它的容量与我们手动启用的goroutine的数量相同，之后再利用这个通道，让主goroutine等待其他 goroutine的运行结束
	// 方式一: 使用channel
	coordinateWithChan()
	fmt.Println()
	// 方式二: 使用sync.WaitGroup
	coordinateWithWaitGroup()
}

func coordinateWithChan() {
	sign := make(chan struct{}, 2)
	num := int32(0)
	fmt.Printf("The number: %d [with chan struct{}]\n", num)
	max := int32(10)
	go addNum(&num, 1, max, func() {
		sign <- struct{}{}
	})
	go addNum(&num, 2, max, func() {
		sign <- struct{}{}
	})
	<-sign
	<-sign
}

func coordinateWithWaitGroup() {
	var wg sync.WaitGroup
	wg.Add(2)
	num := int32(0)
	fmt.Printf("The number: %d [with sync.WaitGroup]\n", num)
	max := int32(10)
	go addNum(&num, 3, max, wg.Done)
	go addNum(&num, 4, max, wg.Done)
	wg.Wait()
}

// addNum 用于原子地增加numP所指的变量的值。
func addNum(numP *int32, id, max int32, deferFunc func()) {
	defer func() {
		deferFunc()
	}()
	for i := 0; ; i++ {
		currNum := atomic.LoadInt32(numP)
		if currNum >= max {
			break
		}
		newNum := currNum + 2
		time.Sleep(time.Millisecond * 200)
		if atomic.CompareAndSwapInt32(numP, currNum, newNum) {
			fmt.Printf("The number: %d [%d-%d]\n", newNum, id, i)
		} else {
			fmt.Printf("The CAS operation failed. [%d-%d]\n", id, i)
		}
	}
}
