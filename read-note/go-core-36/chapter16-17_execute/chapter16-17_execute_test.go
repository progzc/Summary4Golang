package chapter16_17_execute

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestExec_1 程序的执行
// (1)常规概念:
//
//	a.程序: 静静躺着的代码; 进程: 运行着的程序; 线程：线程总是在进程之内的,它可以被视为进程中运行着的控制流
//	b.GMP: G代表用户级线程,M代表系统级线程
//
// (2)关于GMP调度器: G(goroutine,用户级线程)、P(processor,中介)和M(machine,系统级线程)
//
//	a.创建 G 的成本也是非常低的。创建一个 G 并不会像新建一个进程或者一个系统级线程那样，必须通过操作系统的系统调用来完成，
//	  在 Go 语言的运行时系统内部就可以完全做到了，更何况一个 G 仅相当于为需要并发执行代码片段服务的上下文环境而已。
//
// (3)Q:用什么手段可以对 goroutine 的启用数量加以限制?
//
//	A1:一个很简单且很常用的方法是，使用一个通道保存一些令牌。只有先拿到一个令牌，才能启用一个 goroutine。
//	   另外在go函数即将执行结束的时候还需要把令牌及时归还给那个通道。
//	A2:更高级的手段就需要比较完整的设计了。比如，任务分发器 + 任务管道（单层的通道）+ 固定个数的 goroutine。
//	   又比如，动态任务池（多层的通道）+ 动态 goroutine 池（可由前述的那个令牌方案演化而来）。等等
//
// (4)Q:怎么控制P的数量?
//
//	A:runtime.GOMAXPROCS(maxProcs)
//
// (5)Q:runtime包中提供了哪些与模型三要素 G、P 和 M 相关的函数?
//
//	A:https://golang.google.cn/pkg/runtime/
func TestExec_1(t *testing.T) {
	for i := 0; i < 10; i++ {
		// 会异步并发执行
		go func() {
			fmt.Println(i) // 会输出若干个10 (少于10个)
		}()
	}
}

// TestExec_2 怎样才能让主 goroutine 等待其他 goroutine
// 方法一: sleep
func TestExec_2(t *testing.T) {
	num := 10
	for i := 0; i < num; i++ {
		go func() {
			fmt.Println(i) // 会输出10或其他值,数量总共是10个
		}()
	}
	time.Sleep(time.Millisecond * 500)
}

// TestExec_3 怎样才能让主 goroutine 等待其他 goroutine
// 方法二: 使用通道
func TestExec_3(t *testing.T) {
	num := 10
	sign := make(chan struct{}, num)

	for i := 0; i < num; i++ {
		go func() {
			fmt.Println(i)
			sign <- struct{}{}
		}()
	}

	for j := 0; j < num; j++ {
		<-sign
	}
}

// TestExec_4 怎样才能让主 goroutine 等待其他 goroutine
// 方法三: 使用sync.WaitGroup
func TestExec_4(t *testing.T) {
	num := 10
	var w sync.WaitGroup
	w.Add(num)
	for i := 0; i < num; i++ {
		go func() {
			fmt.Println(i)
			w.Done()
		}()
	}
	w.Wait()
}

// TestExec_5 怎样让我们启用的多个 goroutine 按照既定的顺序运行? 即如何使异步发起的go函数得到同步运行
// 方法: 使用 自旋(spinning)+原子操作
func TestExec_5(t *testing.T) {
	// 由于选用的原子操作函数对被操作的数值的类型有约束,所以对count以及相关的变量和参数的类型进行了统一的变更（由int变为了uint32）
	var count uint32
	trigger := func(i uint32, fn func()) {
		for {
			if n := atomic.LoadUint32(&count); n == i {
				fn()
				atomic.AddUint32(&count, 1)
				break
			}
			time.Sleep(time.Nanosecond)
		}
	}
	for i := uint32(0); i < 10; i++ {
		go func(i uint32) {
			// fn 是回调函数
			fn := func() {
				fmt.Println(i)
			}
			// trigger中控制: 当满足条件时,回调fn
			trigger(i, fn)
		}(i)
	}
	// 这一行代码避免主goroutine执行完成后马上退出
	trigger(10, func() {})
}
