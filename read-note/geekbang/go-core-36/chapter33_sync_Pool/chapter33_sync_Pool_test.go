package chapter33_sync_Pool

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"testing"
)

// TestPool_1
// (1)sync.Pool类型（也称为临时对象池）属于结构体类型，它的值在被真正使用之后，就不应该再被复制了(只能传递指针)
// (2)基本注意事项
//
//	a."临时对象"意味着
//		它们的创建和销毁可以在任何时候发生，并且完全不会影响到程序的功能（这决定了sync.Pool的应用场景）;
//		无需被区分,其中的任何一个值都可以代替另一个。
//	b.sync.Pool的Put方法用于在当前的池中存放临时对象，它接受一个interface{}类型的参数。
//	c.sync.Pool的Get方法可能会从当前的池中删除掉任何一个值，然后把这个值作为结果返回（返回的仍然是interface{}类型）;
//	  如果对象池中没有,则会使用New字段代表的函数创建一个新值(这个新值并不会被存入池中)并直接返回。
//	d.若在初始化sync.Pool时没指定New字段，那么调用Get方法时可能会得到nil，故而sync.Pool并非时开箱即用的
//	e.典型例子: 标准库代码包fmt就使用到了sync.Pool类型。
//
// (3)Q:什么说临时对象池中的值会被及时地清理掉?
//
//		  A:Go语言运行时系统中的垃圾回收器，所以在每次开始执行之前，都会对所有已创建的临时对象池中的值进行全面地清除。
//		  对象池中的值的详细的生命周期:
//			a.sync包在被初始化时,会向runtime注册一个函数,这个函数（简称为池清理函数）的功能是清除所有已创建的临时对象池中的值，
//			  注册后,每次进行垃圾回收时,就会执行池清理函数。
//			b.sync包有个包级私有的全局变量(简称为池汇总列表),代表了当前的程序中使用的所有临时对象池的汇总，它是元素类型为*sync.Pool的切片。
//			c.在一个临时对象池的Put方法或Get方法第一次被调用的时候，这个池就会被添加到池汇总列表中。基于此,池清理函数总是能访问到所有正在被真正使用的临时对象池。
//			d.当进行垃圾回收时,池清理函数会遍历池汇总列表。对于其中的每一个临时对象池，它都会先将池中所有的私有临时对象和共享临时对象列表都置为nil，
//		      然后再把这个池中的所有本地池列表都销毁掉。
//			e.池清理函数会把池汇总列表重置为空的切片。如此一来，这些池中存储的临时对象就全部被清除干净了。
//	     总结: 垃圾回收==》池清理函数==》遍历池汇总列表==》讲每个临时对象池中的临时对象置为nil
//
// (4)Q:临时对象池存储值所用的数据结构是怎样的?
//
//		  A:类似于GMP模型，有 本地池列表（其长度等于P的数量,相当于全局队列）、本地池(相当于P队列);
//	     本地池都包含了三个字段(或者说组件)，它们是：存储私有临时对象的字段private、代表了共享临时对象列表的字段shared，以及一个sync.Mutex类型的嵌入字段。
//			一个临时对象池的Put方法或Get方法会获取到哪一个本地池，完全取决于调用它的代码所在的 goroutine 关联的那个 P。
//
// (5)Q:临时对象池是怎样利用内部数据结构来存取值的?
//
//		  A:
//		 	a.临时对象池的Put方法总会先试图把新的临时对象，存储到对应的本地池的private字段中，以便在后面获取临时对象的时候，可以快速地拿到一个可用的值。
//			  只有当这个private字段已经存有某个值时，该方法才会去访问本地池的shared字段。
//			  一个本地池的private字段，只可能被与之对应的那个 P 所关联的 goroutine 中的代码访问到，所以可以说，它是 P 级私有的。
//			b.临时对象池的Get方法，总会先试图从对应的本地池的private字段处获取一个临时对象。只有当这个private字段的值为nil时，它才会去访问本地池的shared字段。
//	       一个本地池的shared字段原则上可以被任何 goroutine 中的代码访问到，不论这个 goroutine 关联的是哪一个 P。
//			c.访问shared字段时，就要加互斥锁sync.Mutex。
//	  举例：Get对象的顺序
//			a.首先从本goroutine关联的对象池中的private字段获取
//			b.第一步获取不到,然后从本goroutine关联的对象池中的shared字段获取（需要加锁）
//			c.第二步获取不到,然后从其他goroutine关联的对象池中的shared字段获取（需要加锁）（work steal）
//			d.第三步获取不到,则调用New字段代表的函数创建并直接返回
//
// (6)Q:怎样保证一个临时对象池中总有比较充足的临时对象?
//
//	A:首先，我们应该事先向临时对象池中放入足够多的临时对象。其次，在用完临时对象之后，我们需要及时地把它归还给临时对象池。
//	  最后，我们应该保证它的New字段所代表的值是可用的。虽然New函数返回的临时对象并不会被放入池中，但是起码能够保证池的Get方法总能返回一个临时对象。
func TestPool_1(t *testing.T) {
	buf := GetBuffer()
	defer buf.Free()
	buf.Write("A Pool is a set of temporary objects that" +
		"may be individually saved and retrieved.")
	buf.Write("A Pool is safe for use by multiple goroutines simultaneously.")
	buf.Write("A Pool must not be copied after first use.")

	fmt.Println("The data blocks in buffer:")
	for {
		block, err := buf.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(fmt.Errorf("unexpected error: %s", err))
		}
		fmt.Print(block)
	}
}

// bufPool 代表存放数据块缓冲区的临时对象池。
var bufPool sync.Pool

// Buffer 代表了一个简易的数据块缓冲区的接口。
type Buffer interface {
	// Delimiter 用于获取数据块之间的定界符。
	Delimiter() byte
	// Write 用于写一个数据块。
	Write(contents string) (err error)
	// Read 用于读一个数据块。
	Read() (contents string, err error)
	// Free 用于释放当前的缓冲区。
	Free()
}

// myBuffer 代表了数据块缓冲区一种实现。
type myBuffer struct {
	buf       bytes.Buffer
	delimiter byte
}

func (b *myBuffer) Delimiter() byte {
	return b.delimiter
}

func (b *myBuffer) Write(contents string) (err error) {
	if _, err = b.buf.WriteString(contents); err != nil {
		return
	}
	return b.buf.WriteByte(b.delimiter)
}

func (b *myBuffer) Read() (contents string, err error) {
	return b.buf.ReadString(b.delimiter)
}

func (b *myBuffer) Free() {
	bufPool.Put(b)
}

// delimiter 代表预定义的定界符。
var delimiter = byte('\n')

func init() {
	bufPool = sync.Pool{
		New: func() interface{} {
			return &myBuffer{delimiter: delimiter}
		},
	}
}

// GetBuffer 用于获取一个数据块缓冲区。
func GetBuffer() Buffer {
	return bufPool.Get().(Buffer)
}
