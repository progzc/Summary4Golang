package chapter13_15_channel

import (
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"testing"
	"time"
)

// 概要
//	(1)黑科技：使用反射操作Channel
//	(2)channel的应用：
//		a.实现消息交流：worker池
//		b.实现数据传递：任务编排之击鼓传花
//		c.实现信号通知(wait/notify功能)、优雅退出
//		d.实现互斥锁
//	(3)Q:在利用chan实现互斥锁的时候，如果buffer设置的不是1，而是一个更大的值，会出现什么状况吗？能解决什么问题吗？
//	   A:chan实现互斥锁，如果buffer大于1，可以实现令牌桶
//	(4)Q:使用chan实现互斥锁，与sync.Mutex的区别？
//	   A:使用chan实现互斥锁，可以很容易实现tryLock、timeout功能；如果不需要这些功能，请优先使用sync.Mutex互斥锁

// TestChannelCodePattern_1
// (1)使用反射操作Channel
func TestChannelCodePattern_1(t *testing.T) {
	var ch1 = make(chan int, 10)
	var ch2 = make(chan int, 10)

	// 创建SelectCase
	var cases = createCases(ch1, ch2)

	// 执行10次select
	for i := 0; i < 10; i++ {
		chosen, recv, ok := reflect.Select(cases)
		if recv.IsValid() { // recv case
			fmt.Println("recv:", cases[chosen].Dir, recv, ok)
		} else { // send case
			fmt.Println("send:", cases[chosen].Dir, ok)
		}
	}
}

// TestChannelCodePattern_2
// (2)worker池：生产者和消费者的消息交流都是通过Channel实现的
//	一个例子：http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/
func TestChannelCodePattern_2(t *testing.T) {
}

// TestChannelCodePattern_3
// (3)任务编排：击鼓传花
// 有4个goroutine，编号为1、2、3、4。每秒钟会有一个goroutine打印出它自己的编号，
// 要求你编写程序，让输出的编号总是按照1、2、3、4、1、2、3、4……这个顺序打印出来。
func TestChannelCodePattern_3(t *testing.T) {
	chs := []chan Token{
		make(chan Token),
		make(chan Token),
		make(chan Token),
		make(chan Token),
	}

	// 创建4个worker
	for i := 0; i < 4; i++ {
		go newWorker(i, chs[i], chs[(i+1)%4])
	}

	//首先把令牌交给第一个worker
	chs[0] <- struct{}{}
	select {}
}

// TestChannelCodePattern_4
// (4)信号通知
//	a.实现wait/notify功能（与sync.Cond的功能一致）
//	b.在程序关闭时，在退出之前做一些清理（doCleanup方法）的动作，即实现程序的graceful shutdown，在退出之前执行一些连接关闭、文件close、
//	  缓存落盘等动作。
func TestChannelCodePattern_4(t *testing.T) {
	var closing = make(chan struct{})
	var closed = make(chan struct{})

	go func() {
		// 模拟业务处理
		for {
			select {
			case <-closing:
				return
			default:
				// ....... 业务计算
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// 处理CTRL+C等中断信号
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	// closing，代表程序退出，但是清理工作还没做
	close(closing)
	// 执行退出之前的清理动作
	go doCleanup(closed)

	select {
	case <-closed:
	case <-time.After(time.Second):
		fmt.Println("清理超时，不等了")
	}
	fmt.Println("优雅退出")
}

// TestChannelCodePattern_5
// (5)互斥锁
//	要想使用chan实现互斥锁，至少有两种方式：
//		a.一种方式是先初始化一个capacity等于1的Channel，然后再放入一个元素。这个元素就代表锁，谁取得了这个元素，就相当于获取了这把锁。
//		b.另一种方式是，先初始化一个capacity等于1的Channel，它的“空槽”代表锁，谁能成功地把元素发送到这个Channel，谁就获取了这把锁。
func TestChannelCodePattern_5(t *testing.T) {
	//a.一种方式是先初始化一个capacity等于1的Channel，然后再放入一个元素。这个元素就代表锁，谁取得了这个元素，就相当于获取了这把锁。
	m := NewMutex()
	ok := m.TryLock()
	fmt.Printf("locked v %v\n", ok)
	ok = m.TryLock()
	fmt.Printf("locked %v\n", ok)
}

// -------------------------------------使用channel实现互斥锁---------------------------------------
// 使用chan实现互斥锁
type Mutex struct {
	ch chan struct{}
}

// 使用锁需要初始化
func NewMutex() *Mutex {
	mu := &Mutex{make(chan struct{}, 1)}
	mu.ch <- struct{}{}
	return mu
}

// 请求锁，直到获取到
func (m *Mutex) Lock() {
	<-m.ch
}

// 解锁
func (m *Mutex) Unlock() {
	select {
	case m.ch <- struct{}{}:
	default:
		panic("unlock of unlocked mutex")
	}
}

// 尝试获取锁
func (m *Mutex) TryLock() bool {
	select {
	case <-m.ch:
		return true
	default:
	}
	return false
}

// 加入一个超时的设置
func (m *Mutex) LockTimeout(timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	select {
	case <-m.ch:
		timer.Stop()
		return true
	case <-timer.C:
	}
	return false
}

// 锁是否已被持有
func (m *Mutex) IsLocked() bool {
	return len(m.ch) == 0
}

// -------------------------------------信号通知---------------------------------------
func doCleanup(closed chan struct{}) {
	time.Sleep(time.Minute)
	// closed，代表清理工作已经做完
	close(closed)
}

// -------------------------------------任务编排---------------------------------------
type Token struct{}

func newWorker(id int, ch chan Token, nextCh chan Token) {
	for {
		token := <-ch       // 取得令牌
		fmt.Println(id + 1) // id从1开始
		time.Sleep(time.Second)
		nextCh <- token
	}
}

// -------------------------------------使用反射操作Channel------------------------------
func createCases(chs ...chan int) []reflect.SelectCase {
	var cases []reflect.SelectCase

	// 创建recv case
	for _, ch := range chs {
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		})
	}

	// 创建send case
	for i, ch := range chs {
		v := reflect.ValueOf(i)
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectSend,
			Chan: reflect.ValueOf(ch),
			Send: v,
		})
	}

	return cases
}
