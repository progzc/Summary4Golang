package chapter08_sync_Once

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"
)

// TestOnce_1
// (1)Q:如何对单例对象进行初始化?
//	  A:有如下几种常见的方式：
//		a.定义package级别的变量,在定义的时候执行初始化。
//		b.在init函数中进行初始化。
//		c.在main函数开始执行的时候，执行一个初始化函数Init。
//		d.使用sync.Once
//		总结：这四种方式都是线程安全的；不同的是，前三种属于立即初始化，最后一种属于延迟初始化。
// (2)实现单例模式的延迟初始化的两种方式：
//	a.使用互斥锁+双重检验
//	  注意事项：
//	  	i)Java中，基于JVM的内存模型，双重检验锁会出现可见性问题,可以通过volatile解决。
//		ii)Go中，使用互斥锁+双重检验实现的单例模式，会出现bug，会有并发安全问题。解决方式是使用 原子操作或原子值atomic.Value。
//	b.使用sync.Once（没有并发安全问题，强烈推荐使用），事实上sync.Once就是Go官方基于安全的双重检验实现的单例模式的体现。
// (3)sync.Once的字段构成
//	a.done：uint32类型，执行原子操作。
//	b.m：sync.Mutex类型，为什么要使用互斥锁?
// 		想一想，看看下面这种方式会带来什么问题?
//			type Once struct {
//				done uint32
//			}
//
//			func (o *Once) Do(f func()) {
//				if !atomic.CompareAndSwapUint32(&o.done, 0, 1) {
//					return
//				}
//				f()
//			}
//			这个实现有一个很大的问题，就是如果参数 f 执行很慢的话，后续调用 Do 方法的 goroutine 虽然看到 done 已经设置为执行过了，
//			但是获取某些初始化资源的时候可能会得到空的资源，因为 f 还没有执行完。所以官方实现在字段中加入了一个sync.Mutex。
// (4)sync.Once的实现原理: 互斥锁+原子操作
//		官方实现如下：
//			func (o *Once) Do(f func()) {
//				// 第一重检查保证性能，快速成功
//				if atomic.LoadUint32(&o.done) == 0 {
//					// Outlined slow-path to allow inlining of the fast-path.
//					o.doSlow(f)
//				}
//			}
//
//			func (o *Once) doSlow(f func()) {
//				o.m.Lock()
//				defer o.m.Unlock()
//				// 第二重检查
//				if o.done == 0 {
//					defer atomic.StoreUint32(&o.done, 1)
//					f()
//				}
//			}
func TestOnce_1(t *testing.T) {
	// 参考文章：https://blog.csdn.net/q5706503/article/details/105870179
	// a.使用 互斥锁+双重检验
	_ = getInstance()
	// b.使用sync.Once
	_ = getInstance2()
}

// TestOnce_2
// (5)sync.Once的注意事项
//	a.sync.Once 只暴露了一个方法 Do，你可以多次调用 Do 方法，但是只有第一次调用 Do 方法时 f 参数才会执行，这里的 f 是一个无参数无返回值的函数。
//	b.死锁：理解了sync.Once的底层使用了sync.Mutex互斥锁就知道这里为什么出现死锁了，应为sync.Mutex不可重复。
//	c.未初始化：如果 f 方法执行的时候 panic，或者 f 执行初始化资源的时候失败了，这个时候，Once 还是会认为初次执行已经成功了，
//			  即使再次调用 Do 方法，也不会再次执行 f。
// (6)sync.Once的使用场景
//	a.Once 常常用来初始化单例资源，或者并发访问只需初始化一次的共享资源，或者在测试的时候初始化一次测试资源。
func TestOnce_2(t *testing.T) {
	//a.sync.Once 只暴露了一个方法 Do，你可以多次调用 Do 方法，
	//但是只有第一次调用 Do 方法时 f 参数才会执行，这里的 f 是一个无参数无返回值的函数。
	var o sync.Once

	// 第一个初始化函数
	f1 := func() {
		fmt.Println("in f1")
	}
	o.Do(f1) // 打印出 in f1

	// 第二个初始化函数
	f2 := func() {
		fmt.Println("in f2")
	}
	o.Do(f2) // 无输出

	// b.死锁：理解了sync.Once的底层使用了sync.Mutex互斥锁就知道这里为什么出现死锁了，应为sync.Mutex不可重复。
	var o2 sync.Once
	o2.Do(func() {
		o2.Do(func() {
			fmt.Println("初始化")
		})
	})
}

// TestOnce_3
// (7)Q:如何避免Once.Do方法可能因为函数执行失败而未初始化资源？
//	a.第一种方法：增加一个err返回参数，为nil标志表示初始化成功，否则初始化失败。
//	b.第二种方法：第一种方法的弊端在于，每次执行Do方法都要check err是否为nil；标准库的 Once 并不会告诉你是否初始化完成了，
//				只是让你放心大胆地去执行 Do 方法。那怎么实现官方的类似操作呢?很简单，增加一个辅助变量，自己去检查是否初始化过了。
//				或者使用hack的方式（即unsafe.Pointer）获取到初始化是否执行成功。
func TestOnce_3(t *testing.T) {
	//a.第一种方法：增加一个err返回参数，为nil标志表示初始化成功，否则初始化失败。
	_ = Once{}
	//b.第二种方法：第一种方法的弊端在于，每次执行Do方法都要check err是否为nil；标准库的 Once 并不会告诉你是否初始化完成了，
	//  只是让你放心大胆地去执行 Do 方法。那怎么实现官方的类似操作呢?很简单，增加一个辅助变量，自己去检查是否初始化过了。
	var o2 Once2
	fmt.Println(o2.IsDone()) //false

	o2.Do(func() {
		time.Sleep(time.Second)
	})
	fmt.Println(o2.IsDone()) //true
}

// ---------------------------如何避免Once.Do方法可能因为函数执行失败而未初始化资源？方法一（有弊端）----------------------
// 一个功能更加强大的Once
type Once struct {
	m    sync.Mutex
	done uint32
}

// 传入的函数f有返回值error，如果初始化失败，需要返回失败的error
// Do方法会把这个error返回给调用者
func (o *Once) Do(f func() error) error {
	if atomic.LoadUint32(&o.done) == 1 { //fast path
		return nil
	}
	return o.slowDo(f)
}

// 如果还没有初始化
func (o *Once) slowDo(f func() error) error {
	o.m.Lock()
	defer o.m.Unlock()
	var err error
	if o.done == 0 { // 双检查，还没有初始化
		err = f()
		if err == nil { // 初始化成功才将标记置为已初始化
			atomic.StoreUint32(&o.done, 1)
		}
	}
	return err
}

// ---------------------------如何避免Once.Do方法可能因为函数执行失败而未初始化资源？方法二（优化版）----------------------
// Once2 是一个扩展的sync.Once类型，提供了一个Done方法
type Once2 struct {
	sync.Once
}

// Done 返回此Once是否执行过
// 如果执行过则返回true
// 如果没有执行过或者正在执行，返回false
func (o *Once2) IsDone() bool {
	return atomic.LoadUint32((*uint32)(unsafe.Pointer(&o.Once))) == 1
}

// -------------------单例模式的实现：使用互斥锁+双重检验------------------------------
type Instance struct{}

var (
	instance *Instance
	lock     sync.Mutex
)

func getInstance() *Instance {
	// 双重检验保证性能
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			i := Instance{}
			instance = &i // Bug: 这一步操作不是原子的, 使用go test -race可以检测出来。解决方式是使用 原子操作或原子值atomic.Value
		}
	}
	return instance
}

// -------------------单例模式的实现：使用sync.Once------------------------------
type Instance2 struct{}

var (
	instance2 *Instance2
	once      sync.Once
)

func getInstance2() *Instance2 {
	once.Do(func() {
		instance2 = &Instance2{}
	})
	return instance2
}
