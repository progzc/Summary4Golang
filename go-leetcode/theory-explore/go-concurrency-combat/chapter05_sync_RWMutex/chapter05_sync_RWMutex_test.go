package chapter05_sync_RWMutex

import (
	"sync"
	"testing"
	"time"
)

// TestSyncRWMutex_1
// (1)readers-writers问题：https://en.wikipedia.org/wiki/Readers%E2%80%93writers_problem
//	概念：同时可能有多个读或者多个写，但是只要有一个线程在执行写操作，其它的线程都不能执行读写操作。
//		 Go标准库中的RWMutex（读写锁）就是用来解决这类readers-writers问题的。
// (2)sync.RWMutex 的基本使用
//	注意事项：RWMutex设计是Write-preferring方案，即一个正在阻塞的Lock调用会排除新的reader请求到锁。
//			即：当writer请求锁的时候，是无法改变既有的reader持有锁的现实的，也不会强制这些reader释放锁，
//			   它的优先权只是限定后来的reader不要和它抢。
//	使用场景：读多写少
// (3)sync.RWMutex 的原理：基于互斥锁实现
//		type RWMutex struct {
//			w           Mutex   // 互斥锁解决多个writer的竞争
//			writerSem   uint32  // writer信号量，作用是阻塞和唤醒writer
//			readerSem   uint32  // reader信号量，作用是阻塞和唤醒reader
//			readerCount int32   // 记录当前reader的数量（以及是否有writer竞争锁(当readerCount小于0时，表示有writer请求锁)）
//			readerWait  int32   // 记录writer请求锁时需要等待read完成的reader的数量
//		}
//		const rwmutexMaxReaders = 1 << 30 // 定义了最大的reader数量
//	总结：
//		1.理解readerCount的两层含义：
//			a.没有writer竞争或持有锁时，readerCount和我们正常理解的reader的计数是一样的。
//			b.但是，如果有writer竞争锁或者持有锁时，那么，readerCount不仅仅承担着reader的计数功能，
//			  还能够标识当前是否有writer竞争或持有锁。
//		2.实现原理：
//			a.RLock：将 readerCount = readerCount+1，然后判断当 readerCount<0 时，说明有writer等待请求锁，将后来的reader阻塞。
//			b.RUnlock：将 readerCount = readerCount-1，然后判断当 readerCount<0 时，说明有writer等待请求锁；
//					   然后调用rUnlockSlow方法，检查是不是reader都释放读锁了（通过readerWait=readerWait-1，readerWait==0判断读锁都释放了），
//					   如果读锁都释放了，那么可以唤醒请求写锁的writer了。
//			c.Lock：先通过readerCount=readerCount-rwmutexMaxReaders 反转readerCount，告诉后来的reader有writer竞争锁；
//					然后判断当前是否有reader持有锁，若有，则阻塞writer，否则加锁完毕。
//			d.Unlock：首先恢复反转，告诉后来的reader没有活跃的writer了；
//					  然后唤醒阻塞的reader，并释放掉自己持有的锁。
//		3.需要注意的是：注意字段的更改和内部互斥锁的顺序关系。
//					  在Lock方法中，是先获取内部互斥锁，才会修改的其他字段；而在Unlock方法中，是先修改的其他字段，才会释放内部互斥锁，
//					  这样才能保证字段的修改也受到互斥锁的保护。

//		// ---------------------------RLock/RUnlock----------------------------------
//		func (rw *RWMutex) RLock() {
//			if atomic.AddInt32(&rw.readerCount, 1) < 0 {
//				// rw.readerCount是负值的时候，意味着此时有writer等待请求锁，因为writer优先级高，所以把后来的reader阻塞休眠
//				runtime_SemacquireMutex(&rw.readerSem, false, 0)
//			}
//		}
//		func (rw *RWMutex) RUnlock() {
//			if r := atomic.AddInt32(&rw.readerCount, -1); r < 0 {
//				rw.rUnlockSlow(r) // 有等待的writer
//			}
//		}
//		func (rw *RWMutex) rUnlockSlow(r int32) {
//			if atomic.AddInt32(&rw.readerWait, -1) == 0 {
//				// 最后一个reader了，writer终于有机会获得锁了
//				runtime_Semrelease(&rw.writerSem, false, 1)
//			}
//		}

//		// ---------------------------Lock/Unlock----------------------------------
//		func (rw *RWMutex) Lock() {
//			// 首先解决其他writer竞争问题
//			rw.w.Lock()
//			// 反转readerCount，告诉reader有writer竞争锁
//			r := atomic.AddInt32(&rw.readerCount, -rwmutexMaxReaders) + rwmutexMaxReaders
//			// 如果当前有reader持有锁，那么需要等待
//			if r != 0 && atomic.AddInt32(&rw.readerWait, r) != 0 {
//				runtime_SemacquireMutex(&rw.writerSem, false, 0)
//			}
//		}
//		func (rw *RWMutex) Unlock() {
//			// 告诉reader没有活跃的writer了
//			r := atomic.AddInt32(&rw.readerCount, rwmutexMaxReaders)
//
//			// 唤醒阻塞的reader们
//			for i := 0; i < int(r); i++ {
//				runtime_Semrelease(&rw.readerSem, false, 0)
//			}
//			// 释放内部的互斥锁
//			rw.w.Unlock()
//		}

// (4)sync.RWMutex的三个易错点
//	a.不可复制
//		解决方案：vet工具
//	b.重入导致死锁，有以下三种情况会导致死锁：
//		i)writer重入调用Lock。
//		ii)在reader的读操作时调用writer的写操作。
//		iii)writer依赖活跃的reader-->活跃的reader依赖新来的reader-->新来的reader依赖writer，从而形成环形依赖。
//	c.释放未加锁的RWMutex
//		解决方案：使用读写锁遵循 "不遗漏不多余" 原则。
func TestSyncRWMutex_1(t *testing.T) {
	var counter Counter
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

// -------------------------------(2)sync.RWMutex 的基本使用----------------------------
// 一个线程安全的计数器
type Counter struct {
	mu    sync.RWMutex
	count uint64
}

// 使用写锁保护
func (c *Counter) Incr() {
	c.mu.Lock()
	c.count++
	c.mu.Unlock()
}

// 使用读锁保护
func (c *Counter) Count() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.count
}
