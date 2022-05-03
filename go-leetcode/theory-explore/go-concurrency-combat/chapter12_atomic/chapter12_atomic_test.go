package chapter12_atomic

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const x int64 = 1 + 1<<33

// TestAtomic_1
// (1)关于原子操作
// 注意事项：有时看起来貌似一个操作是原子操作，但实际上，对于不同的架构来说，情况是不一样的。
func TestAtomic_1(t *testing.T) {
	var i = x
	// 下面这一句使用GOARCH=386的架构去编译，会被拆成两个指令；使用GOARCH=amd64去编译，则是一条指令。
	// 编译：GOARCH=386 go_knowledge tool compile -N -l test.go_knowledge
	// 反编译：GOARCH=386 go_knowledge tool objdump -gnu test.o
	_ = i
}

// TestAtomic_2
// (2)原子包的应用场景
//	相关文章：https://docs.microsoft.com/zh-cn/windows/win32/dxtecharts/lockless-programming
//	a.不涉及到对资源复杂的竞争逻辑
//	b.从配置服务器读取配置并更新
//	c.实现lock-free数据结构（atomic原子操作是实现lock-free数据结构的基石）
// (3)atomic提供的方法
//	atomic为了支持int32、int64、uint32、uint64、uintptr、Pointer（Add方法不支持）类型，
//	分别提供了AddXXX、CompareAndSwapXXX、SwapXXX、LoadXXX、StoreXXX等方法。
//	例如：
//		func AddUint32(addr *uint32, delta uint32) (new uint32)
//		func CompareAndSwapInt32(addr *int32, old, new int32) (swapped bool)
//		func SwapInt32(addr *int32, new int32) (old int32)
//		func LoadInt32(addr *int32) (val int32)
//		func StoreInt32(addr *int32, val int32)
//	注意事项：
//		a.atomic操作的对象是一个地址，你需要把可寻址的变量的地址作为参数传递给方法，而不是把变量的值传递给方法。
//		b.可以利用计算机补码的知识，将减法变为加法 （参加 src/sync/atomic/doc.go_knowledge#AddUint32）
//			例如：针对func AddUint32(addr *uint32, delta uint32) (new uint32)，可以使用 AddUint32(&x, ^uint32(c-1))，
//				 达到x=x-c的效果。尤其是减1这种操作，可以记为 AddUint32(&x, ^uint32(0))
//	Go中的位操作：
//		https://yourbasic.org/golang/bitwise-operator-cheat-sheet/
func TestAtomic_2(t *testing.T) {
	x1 := ^uint32(0)
	// 求负数二进制的方式是根据其对应正数的补码；即将对应正数的二进制，符号位置1，其他位取反加一
	fmt.Printf("x: %b, -1: %b\n", x1, -1) // x: 11111111111111111111111111111111, -1: -1

	// 需要注意的是：c是整数，则有 ^uint32(c - 1) == ^uint32(c) + 1，即 负数的二进制 == 取反(正数二进制)+1 == 取反(正数二进制-1)
	c := 5
	x21 := ^uint32(c - 1)
	fmt.Printf("x21: %b, -5: %b\n", x21, -5) // x21: 11111111111111111111111111111011, -5: -101
	x22 := ^uint32(c) + 1
	fmt.Printf("x22: %b, -5: %b\n", x22, -5) // x22: 11111111111111111111111111111011, -5: -101

	// ^作为二元运算符：异或，包括符号位在内，相同为0，不相同为1。
	// ^作为一元运算符：按位取反，包括符号位在内。
	// x3 = -2的写法
	x3 := ^1
	fmt.Printf("x3: %b\n", x3) // x3: -10
	x4 := 3 &^ 5               // 0011 &^ 0101 = 0010
	fmt.Printf("x4: %d\n", x4) // x4: 2

}

// TestAtomic_3
// (4)Go中的Value类型
//		type Value struct {
//			v interface{}
//		}
//		func (v *Value) Load() (x interface{})
//		func (v *Value) Store(x interface{})
//	常用的场景：配置变更
//	注意事项：不要储存指针，否则会导致内容泄露，引发不安全的操作。
// (5)第三方库的扩展
//	例如：https://github.com/uber-go/atomic
//		a.定义和封装了几种与常见类型相对应的原子操作类型，这些类型提供了原子操作的方法。
//		  这些类型包括Bool、Duration、Error、Float64、Int32、Int64、String、Uint32、Uint64等。
//		b.提供 String、MarshalJSON、UnmarshalJSON 等辅助方法
func TestAtomic_3(t *testing.T) {
	var config atomic.Value
	// 注意事项：不要储存指针，否则会导致内容泄露，引发不安全的操作。
	config.Store(loadNewConfig())
	var cond = sync.NewCond(&sync.Mutex{})

	// 设置新的config
	go func() {
		for {
			time.Sleep(time.Duration(5+rand.Int63n(5)) * time.Second)
			config.Store(loadNewConfig())
			cond.Broadcast() // 通知等待着配置已变更
		}
	}()

	go func() {
		for {
			cond.L.Lock()
			cond.Wait()                 // 等待变更信号
			c := config.Load().(Config) // 读取新的配置
			fmt.Printf("new config: %+v\n", c)
			cond.L.Unlock()
		}
	}()
	select {}
}

// TestAtomic_4
// (6)使用使用atomic实现 Lock-Free queue
// (7)Q:对一个地址的赋值是原子操作吗？
//	  A:这是一个很有趣的问题，如果是原子操作，还要atomic包干什么？官方的文档中并没有特意的介绍，不过，在一些issue或者论坛中，
//	    每当有人谈到这个问题时，总是会被建议用atomic包。
//		a.在现在的系统中，write的地址基本上都是对齐的（aligned）。比如，32位的操作系统、CPU以及编译器，write的地址总是4的倍数，64位的系统总是8的倍数
//		  （还记得WaitGroup针对64位系统和32位系统对state1的字段不同的处理吗）。对齐地址的写，不会导致其他人看到只写了一半的数据，
//		  因为它通过一个指令就可以实现对地址的操作。如果地址不是对齐的话，那么，处理器就需要分成两个指令去处理，如果执行了一个指令，
//		  其它人就会看到更新了一半的错误的数据，这被称做撕裂写（torn write）。
//		  所以，你可以认为赋值操作是一个原子操作，这个“原子操作”可以认为是保证数据的完整性。
//		b.但是，对于现代的多处理多核的系统来说，由于cache、指令重排，可见性等问题，我们对原子操作的意义有了更多的追求。
//		  在多核系统中，一个核对地址的值的更改，在更新到主内存中之前，是在多级缓存中存放的。
//		  这时，多个核看到的数据可能是不一样的，其它的核可能还没有看到更新的数据，还在使用旧的数据。
//		c.多处理器多核心系统为了处理这类问题，使用了一种叫做内存屏障（memory fence 或 memory barrier）的方式。
//		  一个写内存屏障会告诉处理器，必须要等到它管道中的未完成的操作（特别是写操作）都被刷新到内存中，再进行操作。此操作还
//		  会让相关的处理器的CPU缓存失效，以便让它们从主存中拉取最新的值。
//	总结：
//		a.对一个地址的复制因为内存对齐、指令重排、cache、可见性等原因，不是原子操作。
//		b.而atomic包提供的方法会提供内存屏障的功能，所以，atomic不仅仅可以保证赋值的数据完整性，还能保证数据的可见性；
//		  一旦一个核更新了该地址的值，其它处理器总是能读取到它的最新值；
//		  但是，需要注意的是，因为需要处理器之间保证数据的一致性，atomic的操作也是会降低性能的。
func TestAtomic_4(t *testing.T) {
	// (6)使用使用atomic实现 Lock-Free queue
	_ = NewLKQueue()
}

type Config struct {
	NodeName string
	Addr     string
	Count    int32
}

func loadNewConfig() Config {
	return Config{
		NodeName: "北京",
		Addr:     "10.77.95.27",
		Count:    rand.Int31(),
	}
}
