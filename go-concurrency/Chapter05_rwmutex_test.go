package go_concurrency

import (
	"fmt"
	"math"
	"sync"
	"testing"
	"time"
)

// Counter2 一个线程安全的计数器
type Counter2 struct {
	mu    sync.RWMutex
	count uint64
}

// Incr 使用写锁保护
func (c *Counter2) Incr() {
	c.mu.Lock()
	c.count++
	c.mu.Unlock()
}

// Count 使用读锁保护
func (c *Counter2) Count() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.count
}

// TestRWMutex 读写锁的使用
func TestRWMutex(t *testing.T) {
	var counter Counter2
	for i := 0; i < 10; i++ { // 10个reader
		go func() {
			for {
				counter.Count() // 计数器读操作
				time.Sleep(time.Millisecond)
			}
		}()
	}
	for { // 一个writer
		counter.Incr() // 计数器写操作
		time.Sleep(time.Second)
	}
}

func TestMax(t *testing.T) {
	fmt.Println(1 << 30)
	fmt.Println(math.Pow(2, 30))
}
