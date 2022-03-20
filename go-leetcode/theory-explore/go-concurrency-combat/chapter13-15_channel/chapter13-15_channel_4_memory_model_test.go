package chapter13_15_channel

import "testing"

// TestMemoryModel_1
// (1)内存模型
//	背景：CPU指令重排、多级Cache、编译器优化带来的指令重排使得多核访问同一个变量变得非常复杂。
//	概念：并发环境中多 goroutine 读相同变量的时候，变量的可见性条件。
//		 即在什么条件下，goroutine 在读取一个变量的值的时候，能够看到其它 goroutine 对这个变量进行的写的结果。
//	文章：
//		https://www.pengrl.com/p/34119/
//		http://nil.csail.mit.edu/6.824/2016/notes/gomem.pdf
//	目的：
//		a.向广大的程序员提供一种保证，以便他们在做设计和开发程序时，面对同一个数据同时被多个 goroutine 访问的情况，
//		  可以做一些串行化访问的控制，比如使用 Channel 或者 sync 包和 sync/atomic 包中的并发原语。
//		b.允许编译器和硬件对程序做一些优化。这一点其实主要是为编译器开发者提供的保证，这样可以方便他们对 Go 的编译器做优化。
// (2)重排和可见性问题
//	指令重排：指代码并不一定会按照你写的顺序执行。
//	happens-before：定义两个事件（读、写action）的顺序：如果事件e1 happens before事件e2，那么，我们就可以说事件e2在事件e1之后发生（happens after）。
//  			    如果e1不是happens before e2，同时也不happens after e2，那么，我们就可以说事件e1和e2是同时发生的。
// (3)内存模型的特点
//	a.在一个goroutine内部，程序的执行顺序和它们的代码指定的顺序是一样的，即使编译器或者CPU重排了读写顺序，
//	  从行为上来看，也和代码指定的顺序一样。即在单个的goroutine内部，happens-before的关系和代码编写的顺序是一致的。
//	b.对于不同的goroutine来说，重排却会产生非常大的影响。因为Go只保证goroutine内部重排对读写的顺序没有影响。
//	  如果要保证多个goroutine之间对一个共享变量的读写顺序，在Go语言中，可以使用并发原语为读写操作建立happens-before关系，这样就可以保证顺序了。
//	c.在Go语言中，对变量进行零值的初始化就是一个写操作。
//	d.如果对超过机器word（64bit、32bit或者其它）大小的值进行读写，那么，就可以看作是对拆成word大小的几个读写无序进行。
//	e.Go并不提供直接的CPU屏障（CPU fence）来提示编译器或者CPU保证顺序性，而是使用不同架构的内存屏障指令来实现统一的并发原语。
//	f.Go语言中保证的happens-before关系：
//		i)init函数：依赖分析+包变量==》init函数==》main函数
//			1.如果包p导入了包q，那么，q的init函数的执行一定happens before p的任何初始化代码；
//			2.main 函数一定在导入的包的 init 函数之后执行；
//			3.包级别的变量在同一个文件中是按照声明顺序逐个初始化的，除非初始化它的时候依赖其它的变量，
//			  同一个包下的多个文件，会按照文件名的排列顺序进行初始化，这个顺序被定义在Go语言规范中，而不是Go的内存模型规范中。
//		ii)goroutine：启动goroutine的go语句的执行，一定happens before此goroutine内的代码执行。
//		iii)Channel：
//			1.往Channel中的发送操作，happens before从该Channel接收相应数据的动作完成之前。
//			2.close一个Channel的调用，肯定happens before从关闭的Channel中读取出一个零值。
//			3.对于unbuffered的Channel，也就是容量是0的Channel，从此Channel中读取数据的调用一定happens before往此Channel发送数据的调用完成。
//			4.如果Channel的容量是m（m>0），那么，第n个receive一定happens before第n+m个send的完成。
//		iv)Mutex/RWMutex：第n次的m.Unlock一定happens before第n+1次m.Lock方法的返回。
//		v)WaitGroup：Wait方法等到计数值归零之后才返回。
//		vi)Once：对于once.Do(f)调用，f函数的那个单次调用一定happens before任何once.Do(f)调用的返回。
//		vii)atomic：可以保证使用atomic的Load/Store的变量之间的顺序性。
func TestMemoryModel_1_1(t *testing.T) {
	go f() //g1
	g()    //g2
}

func TestMemoryModel_1_2(t *testing.T) {
	go setup()
	// 即使观察到done变成true了，最后读取到的a的值仍然可能为空。
	for !done {
	}
	print(a2)
}

func TestMemoryModel_1_3(t *testing.T) {
	go setup2()
	// 即使main goroutine观察到g不为nil，也可能打印出空的msg
	for g2 == nil {
	}
	print(g2.msg)
}

// -------------------------------指令重排3-----------------------------------
type T struct {
	msg string
}

var g2 *T

func setup2() {
	t := new(T)
	t.msg = "hello, world"
	g2 = t
}

// -------------------------------指令重排2-----------------------------------
var a2 string
var done bool

func setup() {
	a2 = "hello, world"
	done = true
}

// -------------------------------指令重排1-----------------------------------
var a, b int

func f() {
	// 由于指令重排，a和b的赋值先后顺序并不一定是按照代码书写顺序
	a = 1 // w之前的写操作
	b = 2 // 写操作w
}

func g() {
	println(b) // 即使这里打印出2，下面的这行也不一定打印出1（可能打印出0）
	println(a) // 即使上面打印出2，这里也不一定打印出1（可能打印出0）
}
