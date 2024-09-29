package chapter06_sync_WaitGroup

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestWaitGroup_1
// (1)sync.WaitGroup的使用
func TestWaitGroup_1(t *testing.T) {
	var counter Counter

	var wg sync.WaitGroup
	wg.Add(10) // WaitGroup的值设置为10

	for i := 0; i < 10; i++ { // 启动10个goroutine执行加1任务
		go worker(&counter, &wg)
	}
	// 检查点，等待goroutine都完成任务
	wg.Wait()
	// 输出当前计数器的值
	fmt.Println(counter.Count())
}

// TestWaitGroup_2
// (2)sync.WaitGroup的实现原理：计数值+waiter数+信号量+原子操作
//
//	a.字段构成
//		noCopy: 辅助vet工具检查是否存在复制操作。如果你想要自己定义的数据结构不被复制使用，或者说，不能通过 vet 工具检查出复制使用的报警，
//		        就可以通过嵌入 noCopy 这个数据类型来实现。
//				type noCopy struct{}
//				// Lock is a no-op used by -copylocks checker from `go vet`.
//				func (*noCopy) Lock()   {}
//				func (*noCopy) Unlock() {}
//		state1: [3]uint32类型，记录了3个数，分别是计数值、waiter数、信号量
//	b.原理
//		i)Add(delta int)方法：delta可以为负数，使用原子操作将计数值加delta。若计数值加上delta后等于0，还要将waiter数减为0，
//			并通过信号量唤醒阻塞的goroutine。
//		ii)Done()方法：调用Add(-1)
//		iii)Wait()方法：检查计数值是否为0，若是则直接返回；若不是则使用原子操作将waiter数加1，并通过信号量阻塞所在的goroutine。
//	c.sync.WaitGroup是可重用的
//
// (3)使用sync.WaitGroup的常见错误
//
//	a.不能直接复制，只能传递指针。
//	b.计数器设置为负值：一般情况下，有两种方法会导致计数器设置为负数。
//		i)调用Add的时候传递一个负数。如果你能保证当前的计数器加上这个负数后还是大于等于0的话，也没有问题，否则就会导致panic。
//		ii)调用Done方法的次数过多，超过了WaitGroup的计数值。Done调用次数比计数值少会造成死锁，Done调用次数比计数值多会造成panic。
//	c.不期望的Add时机。在使用 WaitGroup 的时候，你一定要遵循的原则就是，等所有的 Add 方法调用之后再调用 Wait，
//	  否则就可能导致 panic 或者不期望的结果。
//	d.前一个 Wait 还没结束就重用 WaitGroup。必须等到上一轮的 Wait 完成之后，才能重用 WaitGroup 执行下一轮的 Add/Wait，
//	  如果你在 Wait 还没执行完的时候就调用下一轮 Add 方法，就有可能出现 panic。
func TestWaitGroup_2(t *testing.T) {

}

// 线程安全的计数器
type Counter struct {
	mu    sync.Mutex
	count uint64
}

// 对计数值加一
func (c *Counter) Incr() {
	c.mu.Lock()
	c.count++
	c.mu.Unlock()
}

// 获取当前的计数值
func (c *Counter) Count() uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

// sleep 1秒，然后计数值加1
func worker(c *Counter, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Second)
	c.Incr()
}
