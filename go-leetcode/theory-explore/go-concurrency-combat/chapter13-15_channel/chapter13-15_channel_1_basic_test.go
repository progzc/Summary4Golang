package chapter13_15_channel

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestChannel_1
// (1)类似于Go的Channel库
//
//	a.https://github.com/docker/libchan
//	b.https://github.com/tylertreat/chan
//
// (2)CSP（Communicating Sequential Process）
//
//	概念：CSP 允许使用进程组件来描述系统，它们独立运行，并且只通过消息传递的方式通信。
//	核心思想：执行业务处理的 goroutine 不要通过共享内存的方式通信，而是要通过 Channel 通信的方式分享数据。
//	论文：https://www.cs.cmu.edu/~crary/819-f09/Hoare78.pdf
//
// (3)Channel的应用场景（重点理解）
//
//	a.数据交流：当作并发的buffer或者queue，解决生产者-消费者问题。多个goroutine可以并发当作生产者（Producer）和消费者（Consumer）。
//	b.数据传递：一个goroutine将数据交给另一个goroutine，相当于把数据的拥有权(引用)托付出去。
//	c.信号通知：一个goroutine可以将信号(closing、closed、data ready等)传递给另一个或者另一组goroutine。
//	d.任务编排：可以让一组goroutine按照一定的顺序并发或者串行的执行，这就是编排的功能。
//	e.锁：利用 Channel 也可以实现互斥锁的机制。
//
// (4)Channel的基本用法
//
//	a.channel的类型：这个箭头总是射向左边的，元素类型总在最右边。如果箭头指向chan，就表示可以往chan中塞数据；如果箭头远离chan，就表示chan会往外吐数据。
//		只能接收：chan string 			// 可以发送接收string
//		只能发送：chan<- struct{} 		// 只能发送struct{}
//		既可以接收又可以发送：<-chan int	// 只能从chan接收int
//	b.通过make，我们可以初始化一个chan，未初始化的chan的零值是nil。容量为0的chan叫做unbuffered chan。
//		make(chan int, 9527)：指定通道的容量为9527
//		make(chan int)：不指定通道的容量，则容量默认为0
//	c.nil是chan的零值，是一种特殊的chan，对值是nil的chan的发送接收调用者总是会阻塞。
//	d.基本操作
//		发送数据：ch <- 2000
//		接收数据：
//			x := <-ch
//			// 第二个值是bool类型，代表是否成功地从chan中读取到一个值。
//			// 如果第二个参数是false，chan已经被close而且chan中没有缓存的数据，这个时候，第一个值是零值。
//			x,ok := <-ch
//		关闭通道：close函数
//		通道的容量：cap函数
//		通道中缓存的还未被取走的元素数量：len函数
//		清空通道：
//			for range ch {
//			}
//		读取通道中的值：
//			for v := range ch {
//				fmt.Println(v)
//			}
//	e.关于赋值：双向通道可以赋值给单向通道，反之则不行。
func TestChannel_1(t *testing.T) {
}

// TestChannel_2
// (5)Channel的实现原理：
//
//	a.数据结构（结构体hchan，主要数据结构是循环队列）：https://github.com/golang/go/blob/master/src/runtime/chan.go#L32
//		字段qcount：uint，循环队列元素的数量。代表chan中已经接收但还没被取走的元素的个数。内建函数len可以返回这个字段的值。
//		字段dataqsiz： uint，循环队列的大小。chan使用一个循环队列来存放元素，循环队列很适合这种生产者-消费者的场景。
//		字段buf：unsafe.Pointer，循环队列的指针。存放元素的循环队列的buffer
//		字段elemsize：uint16，chan中元素的大小。声明的时候就确定了。
//		字段elemtype：*_type，chan中元素的类型。声明的时候就确定了。
//		字段closed：uint32，是否已close
//		字段sendx：uint，send在buf中的索引。处理发送数据的指针在buf中的位置。一旦接收了新的数据，指针就会加上elemsize，移向下一个位置。
//				  buf的总大小是elemsize的整数倍，而且buf是一个循环列表。
//		字段recvx：uint，recv在buf中的索引。处理接收请求时的指针在buf中的位置。一旦取出数据，此指针会移动到下一个位置。
//		字段recvq：waitq，receiver的等待队列。chan是多生产者多消费者的模式，如果消费者因为没有数据可读而被阻塞了，就会被加入到recvq队列中。
//		字段sender：waitq，sender的等待队列。
//		字段lock：mutex，互斥锁
//	总结：存储数据的循环队列+存储发送者的sender等待队列（即生产者队列）+存储接收者的receiver等待队列（即消费者队列）+互斥锁。
//
// (6)channel容易犯的错误
//
//	a.出现panic
//		i)close为nil的chan会panic
//		ii)close已经close的chan会panic
//		iii)send已经close的chan
//	b.goroutine泄漏
//
// (7)关于并发工具的选用原则：
//
//	a.共享资源的并发访问使用传统并发原语。
//	b.复杂的任务编排和消息传递使用 Channel。
//	c.消息通知机制使用Channel，除非只想signal一个goroutine，才使用Cond。
//	d.简单等待所有任务的完成用WaitGroup，也有Channel的推崇者用Channel，都可以。
//	e.需要和Select语句结合，使用Channel。
//	f.需要和超时配合时，使用Channel和Context。
//
// (8)不同通道状态下的操作结果
//
//			nil		empty				full				not full/not empty		closed
//	接收		阻塞		阻塞					read				read					返回未读的元素，读完后返回零值
//	发送		阻塞		write				阻塞					write					panic
//	关闭		panic	closed,无未读元素		close,保留未读元素	close,保留未读元素		panic
func TestChannel_2(t *testing.T) {
	// process会出现goroutine泄露（使用无缓冲的通道很有可能导致内存泄露）
	// 原因在于：
	// 如果发生超时，process函数就返回了，这就会导致unbuffered的chan从来就没有被读取。
	// 我们知道，unbuffered chan必须等reader和writer都准备好了才能交流，否则就会阻塞。
	// 超时导致未读，结果就是子goroutine就阻塞在第7行永远结束不了，进而导致goroutine泄漏。
	// 解决办法：
	// 将unbuffered chan改成容量为1的chan。
	_ = process(time.Second)
}

// TestChannel_3
// 思考题：使用Channel进行任务编排。
// 有四个goroutine，编号为1、2、3、4。每秒钟会有一个goroutine打印出它自己的编号，要求你编写一个程序，
// 让输出的编号总是按照1、2、3、4、1、2、3、4、……的顺序打印出来。
func TestChannel_3(t *testing.T) {
	ch1 := make(chan struct{})
	ch2 := make(chan struct{})
	ch3 := make(chan struct{})
	ch4 := make(chan struct{})
	go func() {
		for {
			fmt.Println("I'm goroutine 1")
			time.Sleep(1 * time.Second)
			ch2 <- struct{}{} // I'm done, you turn
			<-ch1
		}
	}()

	go func() {
		for {
			<-ch2
			fmt.Println("I'm goroutine 2")
			time.Sleep(1 * time.Second)
			ch3 <- struct{}{}
		}

	}()

	go func() {
		for {
			<-ch3
			fmt.Println("I'm goroutine 3")
			time.Sleep(1 * time.Second)
			ch4 <- struct{}{}
		}

	}()

	go func() {
		for {
			<-ch4
			fmt.Println("I'm goroutine 4")
			time.Sleep(1 * time.Second)
			ch1 <- struct{}{}
		}

	}()

	select {}
}

// 内存泄露：一个类似的例子https://github.com/etcd-io/etcd/issues/11256
func process(timeout time.Duration) bool {
	ch := make(chan bool)

	go func() {
		// 模拟处理耗时的业务
		time.Sleep(timeout + time.Second)
		ch <- true // 这里可能会永远block，从而导致内存泄露
		fmt.Println("exit goroutine")
	}()
	select {
	case result := <-ch:
		return result
	case <-time.After(timeout):
		return false
	}
}

func TestChannel(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	chs := make([]chan struct{}, 2)
	for i := 0; i < len(chs); i++ {
		chs[i] = make(chan struct{})
	}
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			select {
			case <-chs[0]:
				fmt.Printf("foo")
				chs[1] <- struct{}{}
			}
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			select {
			case <-chs[1]:
				fmt.Printf("bar\n")
				if i < 9 {
					chs[0] <- struct{}{}
				}
			}
		}
	}()
	chs[0] <- struct{}{}
	wg.Wait()
}
