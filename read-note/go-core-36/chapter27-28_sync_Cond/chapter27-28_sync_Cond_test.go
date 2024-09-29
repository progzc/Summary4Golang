package chapter27_28_sync_Cond

import (
	"log"
	"sync"
	"testing"
	"time"
)

// TestCond_1 使用Signal方法
// (1)条件变量与互斥锁
//
//	 a.条件变量是基于互斥锁的，它必须有互斥锁的支撑才能发挥作用。
//		b.二者作用不同：互斥锁被用来保护临界区和共享资源的;而条件变量用于协调想要访问共享资源的那些线程的
//	   当共享资源的状态发生变化时，条件变量可以被用来通知被互斥锁阻塞的线程。
//
// (2)Q:条件变量怎样与互斥锁配合使用? 条件变量的三个方法：等待通知（wait）、单发通知（signal）和广播通知（broadcast）
//
//	  A:条件变量的初始化离不开互斥锁,并且它的方法有的也是基于互斥锁的
//		补充:
//			a.在利用条件变量等待通知的时候,需要在它基于的那个互斥锁保护下进行
//			b.在进行单发通知或广播通知的时候,需要在对应的互斥锁解锁之后再做这两种操作
//
// (3)sync.Cond并非开箱即用的,且只有NewCond这一种生成方式
// (4)Q:*sync.Cond类型的值可以被传递吗？那sync.Cond类型的值呢?
//
//	  A:*sync.Cond可以,sync.Cond不可以.
//		补充说明:
//		sync.Cond类型的值一旦被使用就不应该再被传递了，传递往往意味着拷贝。 拷贝一个已经被使用过的 sync.Cond 值是很危险的，
//		因为在这份拷贝上调用任何方法都会立即引发 panic。 但是它的指针值是可以被拷贝的。
//
// (5)sync.Cond与channel的比较:
//
//	a.sync.Cond的功能使用channel也可以实现
//	b.区别在于:
//		从用途上来讲,sync.Cond主要用于并发流程上的协同,而chan的主要任务是传递数据.
//		从效率上来讲,sync.Cond是更低层次的工具,效率会更高
//		从易用上来讲,chan在使用上更加方便
//
// (6)Q:条件变量的Wait方法做了什么?
//
//		引申问题:
//			为什么先要锁定条件变量基于的互斥锁，才能调用它的Wait方法?
//				答: 因为条件变量的Wait方法在阻塞当前的 goroutine 之前，会解锁它基于的互斥锁，所以在调用该Wait方法之前，我们必须先锁定那个互斥锁，
//				    否则在调用这个Wait方法时，就会引发一个不可恢复的 panic。
//					如果条件变量的Wait方法不先解锁互斥锁的话，那么就只会造成两种后果：不是当前的程序因 panic 而崩溃，就是相关的 goroutine 全面阻塞。
//			为什么要用for语句来包裹调用其Wait方法的表达式，用if语句不行吗?
//				答: 这主要是为了保险起见。如果一个 goroutine 因收到通知而被唤醒，但却发现共享资源的状态，依然不符合它的要求，
//	    			那么就应该再次调用条件变量的Wait方法，并继续等待下次通知的到来。
//					综上所述，在包裹条件变量的Wait方法的时候，我们总是应该使用for语句。
//	  A:条件变量的Wait方法主要做了以下四件事
//		a.把调用它的 goroutine（也就是当前的 goroutine）加入到当前条件变量的通知队列中。
//		b.解锁当前的条件变量基于的那个互斥锁
//		c.让当前的 goroutine 处于等待状态，等到通知到来时再决定是否唤醒它。此时,这个 goroutine 就会阻塞在调用这个Wait方法的那行代码上。
//		d.如果通知到来并且决定唤醒这个goroutine，那么就在唤醒它之后重新锁定当前条件变量基于的互斥锁。自此之后，当前的goroutine就会继续执行后面的代码。
//
// (7)Q:条件变量的Signal方法和Broadcast方法有哪些异同?
//
//	  A:
//		相同点:
//			a.大致作用相同:条件变量的Signal方法和Broadcast方法都是被用来发送通知的
//		不同点:
//			a.具体作用不同: 条件变量的Signal方法只会唤醒一个因此而等待的goroutine;而条件变量的Broadcast方法却会唤醒所有为此等待的goroutine。
//			b.条件变量的Wait方法总会把当前的 goroutine 添加到通知队列的队尾;Signal方法总会从通知队列的队首开始,
//	    	  被唤醒的 goroutine 一般都是最早等待的那一个。
//			c.要你设置好各个 goroutine 所期望的共享资源状态,使用Broadcast方法总没错。
//			d.与Wait方法不同，条件变量的Signal方法和Broadcast方法并不需要在互斥锁的保护下执行;
//	    	  恰恰相反，我们最好在解锁条件变量基于的那个互斥锁之后，再去调用它的这两个方法
//
// (8)Q:sync.Cond类型中的公开字段L是做什么用的？我们可以在使用条件变量的过程中改变这个字段的值吗？
//
//	A:这个字段代表的是当前的sync.Cond值所持有的那个锁。
//	  我们可以在使用条件变量的过程中改变该字段的值，但是在改变之前一定要搞清楚这样做的影响。
func TestCond_1(t *testing.T) {
	// mailbox 代表信箱。
	// 0代表信箱是空的，1代表信箱是满的。
	var mailbox uint8
	// lock 代表信箱上的锁。
	var lock sync.RWMutex
	// sendCond 代表专用于发信的条件变量。
	sendCond := sync.NewCond(&lock)
	// recvCond 代表专用于收信的条件变量。
	recvCond := sync.NewCond(lock.RLocker())

	// sign 用于传递演示完成的信号。
	sign := make(chan struct{}, 3)
	max := 5
	go func(max int) { // 用于发信。
		defer func() {
			sign <- struct{}{}
		}()
		for i := 1; i <= max; i++ {
			time.Sleep(time.Millisecond * 500)
			lock.Lock()
			for mailbox == 1 {
				sendCond.Wait()
			}
			log.Printf("sender [%d]: the mailbox is empty.", i)
			mailbox = 1
			log.Printf("sender [%d]: the letter has been sent.", i)
			lock.Unlock()
			recvCond.Signal()
		}
	}(max)
	go func(max int) { // 用于收信。
		defer func() {
			sign <- struct{}{}
		}()
		for j := 1; j <= max; j++ {
			time.Sleep(time.Millisecond * 500)
			lock.RLock()
			for mailbox == 0 {
				recvCond.Wait()
			}
			log.Printf("receiver [%d]: the mailbox is full.", j)
			mailbox = 0
			log.Printf("receiver [%d]: the letter has been received.", j)
			lock.RUnlock()
			sendCond.Signal()
		}
	}(max)

	<-sign
	<-sign
}

// TestCond_2 使用Broadcast方法
func TestCond_2(t *testing.T) {
	// mailbox 代表信箱。
	// 0代表信箱是空的，1代表信箱是满的。
	var mailbox uint8
	// lock 代表信箱上的锁。
	var lock sync.Mutex
	// sendCond 代表专用于发信的条件变量。
	sendCond := sync.NewCond(&lock)
	// recvCond 代表专用于收信的条件变量。
	recvCond := sync.NewCond(&lock)

	// send 代表用于发信的函数。
	send := func(id, index int) {
		lock.Lock()
		for mailbox == 1 {
			sendCond.Wait()
		}
		log.Printf("sender [%d-%d]: the mailbox is empty.",
			id, index)
		mailbox = 1
		log.Printf("sender [%d-%d]: the letter has been sent.",
			id, index)
		lock.Unlock()
		recvCond.Broadcast()
	}

	// recv 代表用于收信的函数。
	recv := func(id, index int) {
		lock.Lock()
		for mailbox == 0 {
			recvCond.Wait()
		}
		log.Printf("receiver [%d-%d]: the mailbox is full.",
			id, index)
		mailbox = 0
		log.Printf("receiver [%d-%d]: the letter has been received.",
			id, index)
		lock.Unlock()
		sendCond.Signal() // 确定只会有一个发信的goroutine。
	}

	// sign 用于传递演示完成的信号。
	sign := make(chan struct{}, 3)
	max := 6
	go func(id, max int) { // 用于发信。
		defer func() {
			sign <- struct{}{}
		}()
		for i := 1; i <= max; i++ {
			time.Sleep(time.Millisecond * 500)
			send(id, i)
		}
	}(0, max)
	go func(id, max int) { // 用于收信。
		defer func() {
			sign <- struct{}{}
		}()
		for j := 1; j <= max; j++ {
			time.Sleep(time.Millisecond * 200)
			recv(id, j)
		}
	}(1, max/2)
	go func(id, max int) { // 用于收信。
		defer func() {
			sign <- struct{}{}
		}()
		for k := 1; k <= max; k++ {
			time.Sleep(time.Millisecond * 200)
			recv(id, k)
		}
	}(2, max/2)

	<-sign
	<-sign
	<-sign
}
