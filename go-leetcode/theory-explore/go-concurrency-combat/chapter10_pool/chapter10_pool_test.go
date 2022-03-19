package chapter10_pool

import (
	"bytes"
	"sync"
	"testing"
)

// TestPool_1
// (1)sync.Pool的应用场景
//	a.可以将数据库连接、TCP 的长连接对象进行池化，提升程序的性能
//	  i)buffer池：
//	  	https://github.com/vitessio/vitess/blob/main/go/bucketpool/bucketpool.go
//		https://github.com/valyala/bytebufferpool （提供了校准（calibrate，用来动态调整创建元素的权重）的机制，可以“智能”地调整 Pool 的 defaultSize 和 maxSize）
//		https://github.com/oxtoacart/bpool （最大的特色就是能够保持池子中元素的数量，一旦 Put 的数量多于它的阈值，就会自动丢弃）
//	  ii)http client池：
//		https://pkg.go.dev/net/http#Transport，结构体中的idleConn（但没有使用sync.Pool，而是使用map实现，再配合sync.Mutex使用）
//	  iii)tcp连接池：
//		https://github.com/fatih/pool
//	  iv)数据库连接池：
//		https://github.com/golang/go/blob/4fc3896e7933e31822caa50e024d4e139befc75f/src/database/sql/sql.go#L414
//		https://github.com/bradfitz/gomemcache/blob/master/memcache/memcache.go#L150 (采用 Mutex+Slice 实现 Pool)
//	b.goroutine pool协程池(也是一种常见的并发模式)
//	  i)协程池的常见实现(单个协程2KB，可以扩展到1GB)：https://github.com/golang/go/blob/f296b7a6f045325a230f77e9bda1470b1270f817/src/runtime/proc.go#L120
//		https://github.com/valyala/fasthttp/blob/9f11af296864153ee45341d3f2fe0f5178fd6210/workerpool.go#L16
//		https://github.com/panjf2000/ants
//		https://github.com/alitto/pond (功能很齐全)
//		https://github.com/gammazero/workerpool (可以无限制地提交任务，提供了更便利的 Submit 和 SubmitWait 方法提交任务，还可以提供当前的 worker 数和任务数以及关闭 Pool 的功能。)
//		https://github.com/ivpusic/grpool (创建 Pool 的时候需要提供 Worker 的数量和等待执行的任务的最大数量，任务的提交是直接往 Channel 放入任务)
//		https://github.com/dpaks/goworkers (提供了更便利的 Submit 方法提交任务以及 Worker 数、任务数等查询方法、关闭 Pool 的方法。)
//	注意事项：
//		a.sync.Pool池化的对象池化的对象会在未来的某个时候被毫无预兆地回收掉，这对于数据库长连接等场景是不合适的。
//		  这里可以看一下 TCP连接池、数据库连接池 使用sync.Pool是如何实现且解决垃圾回收的问题。
//		b.sync.Pool 本身就是线程安全的，多个 goroutine 可以并发地调用它的方法存取对象。
//		c.sync.Pool 不可在使用之后再复制使用。
// (2)sync.Pool的基本使用
//	New字段：func() interface{}函数类型。当Get且不存在空闲元素时，就会调用New进行创建（若未设置New,则会返回nil）。
//	Get方法：Get() interface{}。需要注意Get返回的值可能为nil(若未设置New字段)；
//       	当Get且不存在空闲元素时，调用New创建的元素直接返回给Get，不会再放到池中。
//  Put方法：Put(x interface{})。用于将一个元素返还给 Pool，Pool 会把这个元素保存到池中，并且可以复用。
//          但如果 Put 一个 nil 值，Pool 就会忽略这个值。
// (3)sync.Pool的字段构成
//	a.local：unsafe.Pointer类型，主要用来存储空闲的元素。
//		unsafe.Pointer指向的是一个长度为P的元素类型为poolLocalInternal的数组。
//		poolLocalInternal类型包含private和 shared两个字段.
//			private字段：interface{}类型，代表一个缓存的元素，只能被当前P访问，所以不存在并发安全问题。
//			shared字段：结构上是一个local-free的双端队列，可以被任意其他P访问（使用原子操作+uintptr）。
//	b.victim：unsafe.Pointer类型，主要用来存储空闲的元素，相当于垃圾分拣站。
//			  i)每次垃圾回收的时候，Pool 会把 victim 中的对象移除，然后把 local 的数据给 victim，这样的话，local 就会被清空；
//			  ii)victim 就像一个垃圾分拣站。如果victim中的元素如果被Get取走，那么就会活过来；否则在没有被引用时，会被垃圾回收掉
//	 		    （这里的垃圾回收指的是sync.Pool的poolCleanup()方法，并非runtime种的垃圾回收）。
// (4)sync.Pool的实现原理：池的大小等于P的大小（与GMP很像）。
//	Get原理：local当前P上的private获取==>local当前P的shared获取==》local其他P的shared窃取==》
// 	        victim当前P上的private获取==>victim当前P的shared获取==》victim其他P的shared窃取==》New创建
//	Put原理：优先设置local当前P的private，若当前的private有值，则插入到本地队列头部。
// (5)sync.Pool的坑
//	a.内存泄漏：https://github.com/golang/go/issues/23199
// 		在使用 sync.Pool 回收 buffer 的时候，一定要检查回收的对象的大小。如果 buffer 太大，就不要回收了，否则就太浪费了。
//	b.内存浪费：https://github.com/golang/go/blob/617f2c3e35cdc8483b950aa3ef18d92965d63197/src/net/http/server.go#L815-L835
//		要做到物尽其用，尽可能不浪费的话，我们可以将 buffer 池分成几层。
func TestPool_1(t *testing.T) {
	// i)buffer池：https://github.com/gohugoio/hugo/blob/master/bufferpool/bufpool.go (注意这里存在内存泄露的风险)
	// a.内存泄漏
	// 原因：
	// 取出来的 bytes.Buffer 在使用的时候，我们可以往这个元素中增加大量的 byte 数据，这会导致底层的 byte slice 的容量可能会变得很大。
	// 这个时候，即使 Reset 再放回到池子中，这些 byte slice 的容量不会改变，所占的空间依然很大。
	// 而且，因为 Pool 回收的机制，这些大的 Buffer 可能不被回收，而是会一直占用很大的空间，这属于内存泄漏的问题。
	// 解决办法：https://github.com/golang/go/issues/23199
	// 在元素放回时，增加了检查逻辑，改成放回的超过一定大小的 buffer，就直接丢弃掉，不再放到池子中。
	_ = buffers
}

var buffers = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func GetBuffer() *bytes.Buffer {
	return buffers.Get().(*bytes.Buffer)
}

func PutBuffer(buf *bytes.Buffer) {
	buf.Reset()
	buffers.Put(buf)
}
