package go_concurrency

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"
)

// 复制Mutex定义的常量
const (
	mutexLocked      = 1 << iota // 加锁标识位置
	mutexWoken                   // 唤醒标识位置
	mutexStarving                // 锁饥饿标识位置
	mutexWaiterShift = iota      // 标识waiter的起始bit位置
)

// 扩展一个Mutex结构
type Mutex struct {
	sync.Mutex
}

// 尝试获取锁
func (m *Mutex) TryLock() bool {
	// 如果能成功抢到锁
	if atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), 0, mutexLocked) {
		return true
	}
	// 如果处于唤醒、加锁或者饥饿状态，这次请求就不参与竞争了，返回false
	old := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	if old&(mutexLocked|mutexStarving|mutexWoken) != 0 {
		return false
	}
	// 尝试在竞争的状态下请求锁
	new := old | mutexLocked
	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), old, new)
}

// TestTryLock 测试TryLock的实现效果
func TestTryLock(t *testing.T) {
	var mu Mutex
	go func() {
		// 启动一个goroutine持有一段时间的锁
		mu.Lock()
		//time.Sleep(1500 * time.Millisecond) // 休眠1.5s
		time.Sleep(1000 * time.Millisecond) // 休眠1s
		mu.Unlock()
	}()
	time.Sleep(time.Second + time.Millisecond*200) // 休眠1.2s
	ok := mu.TryLock()                             // 尝试获取到锁
	if ok {                                        // 获取成功
		fmt.Println("got the lock")
		// do something
		mu.Unlock()
		return
	}
	// 没有获取到
	fmt.Println("can't get the lock")
}

// -----------------------------------------------------------
// TestGetGoroutineCount 测试获取锁当前持有和等待这把锁的 goroutine 的总数
func TestGetGoroutineCount(t *testing.T) {
	var mu Mutex
	for i := 0; i < 1000; i++ { // 启动1000个goroutine
		go func() {
			mu.Lock()
			time.Sleep(time.Second)
			mu.Unlock()
		}()
	}
	for {
		if mu.Count() == 0 {
			break
		}
		// 输出锁的信息
		fmt.Printf("waitings: %d, isLocked: %t, woken: %t,  starving: %t\n",
			mu.Count(), mu.IsLocked(), mu.IsWoken(), mu.IsStarving())
		time.Sleep(time.Second)
	}
}

// Count 获取当前持有和等待这把锁的 goroutine 的总数
func (m *Mutex) Count() int {
	v := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex))) // 获取state字段的值
	v = v >> mutexWaiterShift                                 //得到等待者的数值
	v = v + (v & mutexLocked)                                 //再加上锁持有者的数量，0或者1
	return int(v)
}

// IsLocked 锁是否被持有
func (m *Mutex) IsLocked() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	return state&mutexLocked == mutexLocked
}

// IsWoken 是否有等待者被唤醒
func (m *Mutex) IsWoken() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	return state&mutexWoken == mutexWoken
}

// IsStarving 锁是否处于饥饿状态
func (m *Mutex) IsStarving() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	return state&mutexStarving == mutexStarving
}

// -----------------------------------------------------------
// TestThreadSafeQueue 测试实现线程安全的队列
func TestThreadSafeQueue(t *testing.T) {
	queue := NewSliceQueue(5)
	queue.Enqueue(1)
	queue.Enqueue("abc")
	fmt.Println(queue.data)
	queue.Dequeue()
	fmt.Println(queue.data)
}

// SliceQueue 一个线程安全的队列结构
type SliceQueue struct {
	data []interface{}
	mu   sync.Mutex
}

// NewSliceQueue 创建一个队列
func NewSliceQueue(n int) (q *SliceQueue) {
	return &SliceQueue{data: make([]interface{}, 0, n)}
}

// Enqueue 把值放在队尾
func (q *SliceQueue) Enqueue(v interface{}) {
	q.mu.Lock()
	q.data = append(q.data, v)
	q.mu.Unlock()
}

// Dequeue 移去队头并返回
func (q *SliceQueue) Dequeue() interface{} {
	q.mu.Lock()
	if len(q.data) == 0 {
		q.mu.Unlock()
		return nil
	}
	v := q.data[0]
	q.data = q.data[1:]
	q.mu.Unlock()
	return v
}
