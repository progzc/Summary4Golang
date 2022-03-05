package chapter29_30_atomic

import (
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"sync/atomic"
	"testing"
	"time"
)

// TestAtomic_1
// (1)互斥锁与原子操作
//	a.互斥锁虽然可以保证临界区中代码的串行执行，但却不能保证这些代码执行的原子性（atomicity）
//	b.在众多的同步工具中，真正能够保证原子性执行的只有原子操作
//	c.从本质上来看,互斥锁是一种悲观锁;而 CAS原子操作+自旋 是一种乐观锁
//	d.使用场景:
//		由于原子操作函数只支持非常有限的数据类型，所以在很多应用场景下，互斥锁往往是更加适合的（或者使用atomic.Value）。
//		不过，一旦我们确定了在某个场景下可以使用原子操作函数，比如：只涉及并发地读写单一的整数类型值，或者多个互不相关的整数类型值，那就不要再考虑互斥锁了。
// (2)原子操作的优缺点:
//	优点:可以完全地消除竞态条件，并能够绝对地保证并发安全性;效率比其他同步工具要快出几个数量级
//	缺点:操作需要足够简单，并且要求快速
// (3)Q:sync/atomic包中提供了几种原子操作？可操作的数据类型又有哪些?
//	  A:
//		sync/atomic包中的函数可以做的原子操作有:加法(add)、比较并交换(compare and swap，简称CAS)、加载(load)、存储(store)和交换(swap)
//		可操作的数据类型: int32、int64、uint32、uint64、uintptr，以及unsafe包中的Pointer(对Pointer类型,未提供进行原子加法操作的函数);
//					    sync/atomic包还提供了一个名为Value的类型，它可以被用来存储任意类型的值。
// (4)Q:atomic.AddInt32函数的第一个参数，对应的一定是那个要被增大的整数。可是，这个参数的类型为什么不是int32而是*int32呢?
//	  A:因为原子操作函数需要的是被操作值的指针，而不是这个值本身；被传入函数的参数值都会被复制，像这种基本类型的值一旦被传入函数，就已经与函数外的那个值毫无关系了。
// (5)Q:用于原子加法操作的函数可以做原子减法吗？比如，atomic.AddInt32函数可以用于减小那个被操作的整数值吗？
//	  A:可以。atomic.AddInt32函数的第二个参数代表差量，它的类型是int32，是有符号的。如果我们想做原子减法，那么把这个差量设置为负整数就可以了。
//	  Q:atomic.AddUint32和atomic.AddUint64函数可以用于减小那个被操作的整数值吗？
//	  A:可以。
//		方法一:先把这个差量转换为有符号的int32类型的值,然后再把该值的类型转换为uint32。即 delta := int32(-3); uint32(delta)
//		方法二:^uint32(-N-1))  (基本原理: 整数在计算机中是以补码的形式存在,结果值的补码相同就意味着表达式的等价)
// (6)Q:比较并交换操作与交换操作相比有什么不同？优势在哪里？
//	  A:比较并交换操作即 CAS 操作，是有条件的交换操作，只有在条件满足的情况下才会进行值的交换。
//      而交换指的是，把新值赋给变量，并返回变量的旧值
func TestAtomic_1(t *testing.T) {
	// 第二个衍生问题的示例。
	//Q:atomic.AddUint32和atomic.AddUint64函数可以用于减小那个被操作的整数值吗？
	//A:可以。
	//	方法一:先把这个差量转换为有符号的int32类型的值,然后再把该值的类型转换为uint32。即 delta := int32(-3); uint32(delta)
	//	方法二:^uint32(-N-1))  (基本原理: 整数在计算机中是以补码的形式存在,结果值的补码相同就意味着表达式的等价)
	num := uint32(18)
	fmt.Printf("The number: %d\n", num) // 18
	delta := int32(-3)
	atomic.AddUint32(&num, uint32(delta))
	fmt.Printf("The number: %d\n", num) // 15
	atomic.AddUint32(&num, ^uint32(-(-3)-1))
	fmt.Printf("The number: %d\n", num) // 12

	fmt.Printf("The two's complement of %d: %b\n",
		delta, uint32(delta)) // -3的补码。
	fmt.Printf("The equivalent: %b\n", ^uint32(-(-3)-1)) // 与-3的补码相同。
	fmt.Println()

	// 第三个衍生问题的示例。
	// (6)Q:比较并交换操作与交换操作相比有什么不同？优势在哪里？
	//	  A:比较并交换操作即 CAS 操作，是有条件的交换操作，只有在条件满足的情况下才会进行值的交换。
	//      而交换指的是，把新值赋给变量，并返回变量的旧值
	forAndCAS1()
	fmt.Println()
	forAndCAS2()
}

// forAndCAS1 用于展示简易的自旋锁。
func forAndCAS1() {
	sign := make(chan struct{}, 2)
	num := int32(0)
	fmt.Printf("The number: %d\n", num)
	go func() { // 定时增加num的值。
		defer func() {
			sign <- struct{}{}
		}()
		for {
			time.Sleep(time.Millisecond * 500)
			newNum := atomic.AddInt32(&num, 2)
			fmt.Printf("The number: %d\n", newNum)
			if newNum == 10 {
				break
			}
		}
	}()
	go func() { // 定时检查num的值，如果等于10就将其归零。
		defer func() {
			sign <- struct{}{}
		}()
		// for + CAS == 自旋
		for {
			if atomic.CompareAndSwapInt32(&num, 10, 0) {
				fmt.Println("The number has gone to zero.")
				break
			}
			time.Sleep(time.Millisecond * 500)
		}
	}()
	<-sign
	<-sign
}

// forAndCAS2 用于展示一种简易的（且更加宽松的）互斥锁的模拟。
func forAndCAS2() {
	sign := make(chan struct{}, 2)
	num := int32(0)
	fmt.Printf("The number: %d\n", num)
	max := int32(20)
	go func(id int, max int32) { // 定时增加num的值。
		defer func() {
			sign <- struct{}{}
		}()
		for i := 0; ; i++ {
			currNum := atomic.LoadInt32(&num)
			if currNum >= max {
				break
			}
			newNum := currNum + 2
			time.Sleep(time.Millisecond * 200)
			if atomic.CompareAndSwapInt32(&num, currNum, newNum) {
				fmt.Printf("The number: %d [%d-%d]\n", newNum, id, i)
			} else {
				fmt.Printf("The CAS operation failed. [%d-%d]\n", id, i)
			}
		}
	}(1, max)
	go func(id int, max int32) { // 定时增加num的值。
		defer func() {
			sign <- struct{}{}
		}()
		for j := 0; ; j++ {
			currNum := atomic.LoadInt32(&num)
			if currNum >= max {
				break
			}
			newNum := currNum + 2
			time.Sleep(time.Millisecond * 200)
			if atomic.CompareAndSwapInt32(&num, currNum, newNum) {
				fmt.Printf("The number: %d [%d-%d]\n", newNum, id, j)
			} else {
				fmt.Printf("The CAS operation failed. [%d-%d]\n", id, j)
			}
		}
	}(2, max)
	<-sign
	<-sign
}

// TestAtomic_2
// (1)Q:假设我已经保证了对一个变量的写操作都是原子操作，比如：加或减、存储、交换等等，那我对它进行读操作的时候，还有必要使用原子操作吗？
//    A:很有必要。其中的道理你可以对照一下读写锁。为什么在读写锁保护下的写操作和读操作之间是互斥的？这是为了防止读操作读到没有被修改完的值，对吗？
// (2)Q:怎样用好sync/atomic.Value?
//    A:atomic.Value类型是开箱即用，且只有两个指针方法:Store和Load。
//	  注意事项:
//		a.一旦atomic.Value类型的值（以下简称原子值）被真正使用(即已经用来存储值了)，它就不应该再被复制了。
//		b.不能用原子值存储nil。即不能把nil作为参数值传入原子值的Store方法，否则就会引发一个 panic。
//			特殊地,如果有一个接口类型的变量，它的动态值是nil，但动态类型却不是nil，那么它的值就不等于nil;这样一个变量的值是可以被存入原子值的.
//		c.我们向原子值存储的第一个值，决定了它今后能且只能存储哪一个类型的值;否则在调用Store方法时会引发一个panic.
//			特殊地,先存储一个接口类型的值，然后再存储这个接口的某个实现类型的值,这是不允许的（因为原子值内部是依据被存储值的实际类型来做判断）。
//		    即使是实现了同一个接口的不同类型，它们的值也不能被先后存储到同一个原子值中。
//    几点建议:
//		a.不要把内部使用的原子值暴露给外界,原子值最好是包级私有。
//		b.或者可以声明一个包级私有的原子变量，然后再通过一个或多个公开的函数，让外界间接地使用到它（但不要传递其指针值）
//		c.如果通过某个函数可以向内部的原子值存储值的话，那么就应该在这个函数中先判断被存储值类型的合法性。若不合法，则应该直接返回对应的错误值，从而避免 panic 的发生。
//		d.可能的话，可以把原子值封装到一个数据类型中，比如一个结构体类型。这样，既可以通过该类型的方法更加安全地存储值，又可以在该类型中包含可存储值的合法类型信息。
//		e.尽量不要向原子值中存储引用类型的值。因为这很容易造成安全漏洞。
// (3)Q:如果要对原子值和互斥锁进行二选一，你认为最重要的三个决策条件应该是什么？
//    A:在搞清楚下述问题（以及你关注的其他问题）之后，优先使用原子值。
//		a.被保护的数据是什么类型的？是值类型的还是引用类型的？ // 值类型优先选原子值
//		b.操作被保护数据的方式是怎样的？是简单的读和写还是更复杂的操作？ // 简单的读和写优先选用原子值
//		c.操作被保护数据的代码是集中的还是分散的？如果是分散的，是否可以变为集中的？ // 分散的可以优先选用原子值
func TestAtomic_2(t *testing.T) {
	// 示例1。
	var box atomic.Value
	fmt.Println("Copy box to box2.")
	box2 := box // 原子值在真正使用前可以被复制。
	v1 := [...]int{1, 2, 3}
	fmt.Printf("Store %v to box.\n", v1)
	box.Store(v1)
	fmt.Printf("The value load from box is %v.\n", box.Load())   // [1 2 3]
	fmt.Printf("The value load from box2 is %v.\n", box2.Load()) // nil
	fmt.Println()

	// 示例2。
	v2 := "123"
	fmt.Printf("Store %q to box2.\n", v2)
	box2.Store(v2)                                               // 这里并不会引发panic。
	fmt.Printf("The value load from box is %v.\n", box.Load())   // [1 2 3]
	fmt.Printf("The value load from box2 is %q.\n", box2.Load()) // 123
	fmt.Println()

	// 示例3。
	fmt.Println("Copy box to box3.")
	box3 := box                                                  // 原子值在真正使用后不应该被复制！
	fmt.Printf("The value load from box3 is %v.\n", box3.Load()) // [1 2 3]
	v3 := 123
	fmt.Printf("Store %d to box3.\n", v3)
	//box3.Store(v3) // 这里会引发一个panic，报告存储值的类型不一致。
	_ = box3
	fmt.Println()

	// 示例4。
	var box4 atomic.Value
	v4 := errors.New("something wrong")
	fmt.Printf("Store an error with message %q to box4.\n", v4)
	box4.Store(v4)
	v41 := io.EOF
	fmt.Println("Store a value of the same type to box4.")
	box4.Store(v41)
	v42, ok := interface{}(&os.PathError{}).(error)
	if ok {
		fmt.Printf("Store a value of type %T that implements error interface to box4.\n", v42)
		//box4.Store(v42) // 这里会引发一个panic，报告存储值的类型不一致。
	}
	fmt.Println()

	// 示例5。
	box5, err := NewAtomicValue(v4)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
	fmt.Printf("The legal type in box5 is %s.\n", box5.TypeOfValue())
	fmt.Println("Store a value of the same type to box5.")
	err = box5.Store(v41)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
	fmt.Printf("Store a value of type %T that implements error interface to box5.\n", v42)
	err = box5.Store(v42)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
	fmt.Println()

	// 示例6。
	// e.尽量不要向原子值中存储引用类型的值。因为这很容易造成安全漏洞。
	var box6 atomic.Value
	v6 := []int{1, 2, 3}
	fmt.Printf("Store %v to box6.\n", v6)
	box6.Store(v6)
	v6[1] = 4 // 注意，此处的操作不是并发安全的！
	fmt.Printf("The value load from box6 is %v.\n", box6.Load())
	// 正确的做法如下。
	v6 = []int{1, 2, 3}
	store := func(v []int) {
		replica := make([]int, len(v))
		copy(replica, v)
		box6.Store(replica)
	}
	fmt.Printf("Store %v to box6.\n", v6)
	store(v6)
	v6[2] = 5 // 此处的操作是安全的。
	fmt.Printf("The value load from box6 is %v.\n", box6.Load())
}

type atomicValue struct {
	v atomic.Value
	t reflect.Type
}

func NewAtomicValue(example interface{}) (*atomicValue, error) {
	if example == nil {
		return nil, errors.New("atomic value: nil example")
	}
	return &atomicValue{
		t: reflect.TypeOf(example),
	}, nil
}

func (av *atomicValue) Store(v interface{}) error {
	if v == nil {
		return errors.New("atomic value: nil value")
	}
	t := reflect.TypeOf(v)
	if t != av.t {
		return fmt.Errorf("atomic value: wrong type: %s", t)
	}
	av.v.Store(v)
	return nil
}

func (av *atomicValue) Load() interface{} {
	return av.v.Load()
}

func (av *atomicValue) TypeOfValue() reflect.Type {
	return av.t
}
