package chapter01_04_sync_Mutex

import (
	"fmt"
	"github.com/petermattis/goid"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestMutex_1
// (1)同步原语（互斥锁Mutex、读写锁RWMutex、并发编排WaitGroup、条件变量Cond、Channel）的适用场景
//	a.共享资源：并发地读写共享资源，会出现数据竞争（data race）的问题，所以需要Mutex、RWMutex这样的并发原语来保护。
//	b.任务编排：需要goroutine按照一定的规律执行，而goroutine之间有相互等待或者依赖的顺序关系，我们常常使用WaitGroup或者Channel来实现。
//	c.消息传递：信息交流以及不同的goroutine之间的线程安全的数据交流，常常使用Channel来实现。
// (2)sync.Mutex的基本使用方法
//	举例：count++不是一个原子操作。它至少包含几个步骤，比如读取变量count的当前值，对这个值加1，把结果再保存到count中。
//		 因为不是原子操作，就可能有并发的问题。
//	一个检测并发访问共享资源是否有问题的工具：https://go.dev/blog/race-detector，使用方法如下：
//		go_knowledge run -race counter.go_knowledge
//		go_knowledge test -race counter.go_knowledge
//		go_knowledge compile -race counter.go_knowledge (由于这个工具实现方式，只能通过真正对实际地址进行读写访问的时候才能探测，所以它并不能在编译的时候发现data race的问题)。
//		go_knowledge tool compile -race -S counter.go_knowledge // 将汇编列表打印到标准输出
//			go_knowledge tool compile工具的使用：https://www.cnblogs.com/linguoguo/p/11699006.html
// (3)Q:如果Mutex已经被一个goroutine获取了锁，其它等待中的goroutine们只能一直等待。
//		那么，等这个锁释放后，等待中的goroutine中哪一个会优先获取Mutex呢？
//	  A:等待的goroutine们是以FIFO排队的。详见https://github.com/golang/go_knowledge/blob/master/src/sync/mutex.go_knowledge#L42
func TestMutex_1(t *testing.T) {
	// 互斥锁保护计数器
	var mu sync.Mutex
	// 计数器的值
	var count = 0

	// 辅助变量，用来确认所有的goroutine都完成
	var wg sync.WaitGroup
	wg.Add(10)

	// 启动10个gourontine
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			// 累加10万次
			for j := 0; j < 100000; j++ {
				mu.Lock()
				count++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	fmt.Println(count)
}

// TestMutex_2
// (3)sync.Mutex的是实现原理
//	sync.Mutex设计发展的四个阶段：
//		a.初版：使用一个key字段标记是否持有锁。（利用cas原子操作+信号量实现）
//			type Mutex struct {
//				key  int32		// 持有和等待锁的数量。0：锁未被持有；1：锁被持有，没有等待者；n：锁被持有，还有n-1个等待者
//				sema int32		// 等待者队列使用的信号量，用以阻塞/唤醒goroutine
//			}
//			问题1：Unlock方法可以被任意的goroutine调用释放锁，即使是没持有这个互斥锁的goroutine，也可以进行这个操作。
//				 这是因为，Mutex本身并没有包含持有这把锁的goroutine的信息，所以，Unlock也不会对此进行检查。
//				 解决方案：严格遵循"谁申请，谁释放"的原则。
//			问题2：请求锁的goroutine会排队（FIFO）等待获取互斥锁，看似公平；但是若能把锁交给正在占用CPU时间片的goroutine的话，就不需要做
//				  上下文切换，能够带来更好的性能。

//		b.给新人机会：新的goroutine也能有机会竞争锁。（利用cas原子操作+信号量实现）
//			type Mutex struct {
//				state int32		// 32位从左到右：30位标识阻塞等待的waiter数量、1位标识是否有唤醒的goroutine、1位标识锁是否被持有。
//				sema  uint32	// 等待者队列使用的信号量，用以阻塞/唤醒goroutine
//			}
//			const (
//				mutexLocked = 1 << iota 	// 标识state字段右边第1位代表 锁是否被持有。
//				mutexWoken					// 标识state字段右边第2位代表 是否有唤醒的goroutine。
//				mutexWaiterShift = iota		// 标识state字段右边第3位~32位代表 等待此锁的goroutine数。
//			)

//		c.多给些机会：新来的和被唤醒的有更多的机会竞争锁。（利用cas原子操作+信号量实现+自旋）
//					通过增加自旋，新来的和被唤醒的有更多机会获取到锁，提高了性能（这是因为临界区的代码大都耗时很短，锁可以很快释放，这样只需
//					浪费少许CPU资源就能避免昂贵的上下文切换，因此提升了程序性能）。
//			问题1：有这样一种场景，唤醒的waiter和新年的goroutine竞争锁，每次都失败了，就会产生锁饥饿的问题。

//		d.解决饥饿：解决竞争问题，不会让goroutine长久等待。（利用cas原子操作+信号量实现）
//			type Mutex struct {
//				state int32
//				sema  uint32
//			}
//			const (
//				mutexLocked = 1 << iota 		// 标识state字段右边第1位代表 锁是否被持有。
//				mutexWoken						// 标识state字段右边第2位代表 是否有唤醒的goroutine。
//				mutexStarving					// 标识state字段右边第3位代表 是否有饥饿的goroutine。
//				mutexWaiterShift      = iota	// 标识state字段右边第4位~32位代表 等待此锁的goroutine数。
//				starvationThresholdNs = 1e6		// 饥饿的阈值是1ms
//			)

//	总结，当前sync.Mutex的原理如下：（大致描述出来下面的流程即可）
//		1.Mutex中含有两个字段：
//			state字段：从右到左分别是，右边第1位代表锁是否被持有、右边第2位代表是否有唤醒的goroutine、右边第3位代表是否有饥饿的goroutine
//					  剩下代表等待此锁的goroutine数。
//			sema字段：等待者队列使用的信号量，用以阻塞/唤醒goroutine。
//		2.获取锁的过程：具体流程可详见https://github.com/golang/go_knowledge/blob/master/src/sync/mutex.go_knowledge#L42
//			a.若当前没有goroutine持有锁、没有唤醒的goroutine、没有锁饥饿、没有等待持有锁的goroutine，那么新来的goroutine就会直接获取到锁。
//			b.若有唤醒的goroutine，则唤醒的goroutine与新来的goroutine争抢锁，争抢的方式采用自旋（这是因为临界区的代码大都耗时很短，锁可以很快释放，这样只需
//		  	  浪费少许CPU资源就能避免昂贵的上下文切换，因此提升了程序性能）。
//			c.若新来的goroutine争抢锁失败，那就直接放到等待队列（FIFO）的队尾。
//			d.若唤醒的goroutine一直争抢锁失败，达到饥饿的阈值1ms时，那么就将饥饿标志位置1，唤醒的goroutine锁直接竞争成功，获取到锁，而新来的goroutine
//			  直接放到等待队列（FIFO）的队尾。
//			g.持有锁的goroutine释放锁时，若没有waiter，则直接返回；若有waiter，但没有唤醒的waiter，则唤醒一个waiter。
// (4)Q:等待一个Mutex的goroutine数最大是多少？是否能满足现实的需求？
//	  A:单从程序来看，可以支持 1<<(32-3)-1（其中32为state的类型int32，3位waiter字段的shift），约 0.5 Billion个。
//		考虑到实际goroutine初始化的空间为2K，0.5 Billion * 2K 达到了1TB，单从内存空间来说已经要求极高了，当前的设计肯定可以满足了。
func TestMutex_2(t *testing.T) {
}

// TestMutex_3
// (5)常见的4种错误场景
//	a.Lock/Unlock不是成对出现，具体而言，有以下几种可能：
//		i)代码中有太多的if-else分支，可能在某个分支中漏写了Unlock
//		ii)在重构的时候把Unlock给删除了
//		iii)Unlock误写成了Lock
//	b.Copy已使用的Mutex。
//	  原因在于，Mutex是一个有状态的对象，它的state字段记录这个锁的状态。如果你要复制一个已经加锁的Mutex给一个新的变量，
//	  那么新的刚初始化的变量居然被加锁了，这显然不符合你的期望，因为你期望的是一个零值的Mutex。
//	  这种错误使用方式很可能造成死锁。
//	  解决方式：
//		i)编译时，使用共go vet copy.go检查即可。（推荐这种方式，最好不要依赖于运行时检查机制）
//		  机制原理：https://github.com/golang/tools/blob/master/go/analysis/passes/copylock/copylock.go
//	  	ii)Go在运行时，有死锁检查机制，能够发现死锁的goroutine；
//	  	  详见https://github.com/golang/go_knowledge/blob/master/src/runtime/proc.go_knowledge#L4935
//	c.sync.Mutex不可重入，否则会导致死锁。
//	  当然，我们通过hack设计一个可重入锁，这里有两种方式：
//	  	i)方案一：通过hacker的方式获取到goroutine id，记录下获取锁的goroutine id，它可以实现Locker接口。
//				 一个支持获取多个Go版本的goroutine id的开源库：https://github.com/petermattis/goid
//		ii)方案二：调用Lock/Unlock方法时，由goroutine提供一个token，用来标识它自己，而不是我们通过hacker的方式获取到goroutine id，
//				  但是，这样一来，就不满足Locker接口了。
//	d.死锁：四个必要条件：互斥 + 持有和等待 + 不可剥夺 + 环路等待
func TestMutex_3(t *testing.T) {
	// 设计一个可重入锁

	// i)方案一：通过hacker的方式获取到goroutine id，记录下获取锁的goroutine id，它可以实现Locker接口。
	// 一个支持获取多个Go版本的goroutine id的开源库：https://github.com/petermattis/goid
	_ = RecursiveMutex{}

	// ii)方案二：调用Lock/Unlock方法时，由goroutine提供一个token，用来标识它自己，而不是我们通过hacker的方式获取到goroutine id，
	// 但是，这样一来，就不满足Locker接口了。
	_ = TokenRecursiveMutex{}
}

// TestMutex_4 演示死锁
func TestMutex_4(t *testing.T) {
	// 派出所证明
	var psCertificate sync.Mutex
	// 物业证明
	var propertyCertificate sync.Mutex

	var wg sync.WaitGroup
	wg.Add(2) // 需要派出所和物业都处理

	// 派出所处理goroutine
	go func() {
		defer wg.Done() // 派出所处理完成

		psCertificate.Lock()
		defer psCertificate.Unlock()

		// 检查材料
		time.Sleep(5 * time.Second)
		// 请求物业的证明
		propertyCertificate.Lock()
		propertyCertificate.Unlock()
	}()

	// 物业处理goroutine
	go func() {
		defer wg.Done() // 物业处理完成

		propertyCertificate.Lock()
		defer propertyCertificate.Unlock()

		// 检查材料
		time.Sleep(5 * time.Second)
		// 请求派出所的证明
		psCertificate.Lock()
		psCertificate.Unlock()
	}()

	wg.Wait()
	fmt.Println("成功完成")
}

// ------------------------------------设计可重入锁（方案一）-------------------------------------------------------------
// RecursiveMutex 包装一个Mutex,实现可重入
type RecursiveMutex struct {
	sync.Mutex
	owner     int64 // 当前持有锁的goroutine id
	recursion int32 // 这个goroutine 重入的次数
}

func (m *RecursiveMutex) Lock() {
	gid := goid.Get()
	// 如果当前持有锁的goroutine就是这次调用的goroutine,说明是重入
	if atomic.LoadInt64(&m.owner) == gid {
		m.recursion++
		return
	}
	m.Mutex.Lock()
	// 获得锁的goroutine第一次调用，记录下它的goroutine id,调用次数加1
	atomic.StoreInt64(&m.owner, gid)
	m.recursion = 1
}

func (m *RecursiveMutex) Unlock() {
	gid := goid.Get()
	// 非持有锁的goroutine尝试释放锁，错误的使用
	if atomic.LoadInt64(&m.owner) != gid {
		panic(fmt.Sprintf("wrong the owner(%d): %d!", m.owner, gid))
	}
	// 调用次数减1
	m.recursion--
	if m.recursion != 0 { // 如果这个goroutine还没有完全释放，则直接返回
		return
	}
	// 此goroutine最后一次调用，需要释放锁
	atomic.StoreInt64(&m.owner, -1)
	m.Mutex.Unlock()
}

// ------------------------------------设计可重入锁（方案二）-------------------------------------------------------------
// Token方式的递归锁
type TokenRecursiveMutex struct {
	sync.Mutex
	token     int64
	recursion int32
}

// 请求锁，需要传入token
func (m *TokenRecursiveMutex) Lock(token int64) {
	if atomic.LoadInt64(&m.token) == token { //如果传入的token和持有锁的token一致，说明是递归调用
		m.recursion++
		return
	}
	m.Mutex.Lock() // 传入的token不一致，说明不是递归调用
	// 抢到锁之后记录这个token
	atomic.StoreInt64(&m.token, token)
	m.recursion = 1
}

// 释放锁
func (m *TokenRecursiveMutex) Unlock(token int64) {
	if atomic.LoadInt64(&m.token) != token { // 释放其它token持有的锁
		panic(fmt.Sprintf("wrong the owner(%d): %d!", m.token, token))
	}
	m.recursion--         // 当前持有这个锁的token释放锁
	if m.recursion != 0 { // 还没有回退到最初的递归调用
		return
	}
	atomic.StoreInt64(&m.token, 0) // 没有递归调用了，释放锁
	m.Mutex.Unlock()
}
