package chapter26_sync_Mutex

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"sync"
	"testing"
)

// TestMutex_1
// (1)几个概念
//
//	a.竞态条件: 一旦数据被多个线程共享，那么就很可能会产生争用和冲突的情况
//	b.共享数据的一致性: 多个线程对共享数据的操作总是可以达到它们各自预期的效果
//	c.同步的作用:
//		一是避免是避免多个线程在同一时刻操作同一个数据块;
//		二是协调多个线程,以避免它们在同一时刻执行同一个代码块.
//	d.临界区: 只要一个代码片段需要实现对共享资源的串行化访问
//
// (2)使用互斥锁时有哪些注意事项
//
//		a.不要重复锁定互斥锁 (不可重入)
//		b.不要忘记解锁互斥锁,必要时使用defer语句 (可能会导致死锁或程序崩溃)
//		c.不要对尚未锁定或者已解锁的互斥锁解锁 (这回立即引发panic)
//		d.不要在多个函数之间直接传递互斥锁 (因为sync.Mutex是结构体,属于值类型;原值和它的副本,以及多个副本之间都是完全独立的，它们都是不同的互斥锁)
//	 e.避免死锁最简单的方式是:让每一个互斥锁都只保护一个临界区或一组相关临界区
//
// (3)sync.Mutex 是开箱即用的,我们可以借鉴这种设计方式
// (4)Q:读写锁与互斥锁有哪些异同?
//
//	  A:
//	 	a.sync.RWMutex的写锁会导致读/写锁阻塞,而读锁不会导致读锁阻塞;sync.Mutex都是互斥的
//		b.sync.RWMutex是对sync.Mutex加锁粒度的细化
//
// (5)Q:你知道互斥锁和读写锁的指针类型都实现了哪一个接口吗?
//
//	A:都实现了sync.Locker接口
//
// (6)Q:怎样获取读写锁中的读锁?
//
//	A:sync.RWMutex类型有一个名为RLocker的指针方法可以获取其读锁
func TestMutex_1(t *testing.T) {
	// buffer 代表缓冲区。
	var (
		buffer     bytes.Buffer
		protecting = 1
	)

	const (
		max1 = 5  // 代表启用的goroutine的数量。
		max2 = 10 // 代表每个goroutine需要写入的数据块的数量。
		max3 = 10 // 代表每个数据块中需要有多少个重复的数字。
	)

	// mu 代表以下流程要使用的互斥锁。
	var mu sync.Mutex
	// sign 代表信号的通道。
	sign := make(chan struct{}, max1)

	for i := 1; i <= max1; i++ {
		go func(id int, writer io.Writer) {
			defer func() {
				sign <- struct{}{}
			}()
			for j := 1; j <= max2; j++ {
				// 准备数据。
				header := fmt.Sprintf("\n[id: %d, iteration: %d]",
					id, j)
				data := fmt.Sprintf(" %d", id*j)
				// 写入数据。
				if protecting > 0 {
					mu.Lock()
				}
				_, err := writer.Write([]byte(header))
				if err != nil {
					log.Printf("error: %s [%d]", err, id)
				}
				for k := 0; k < max3; k++ {
					_, err := writer.Write([]byte(data))
					if err != nil {
						log.Printf("error: %s [%d]", err, id)
					}
				}
				if protecting > 0 {
					mu.Unlock()
				}
			}
		}(i, &buffer)
	}

	for i := 0; i < max1; i++ {
		<-sign
	}
	data, err := ioutil.ReadAll(&buffer)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	log.Printf("The contents:\n%s", data)
}
