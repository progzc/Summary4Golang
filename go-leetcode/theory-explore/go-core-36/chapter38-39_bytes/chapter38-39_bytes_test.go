package chapter38_39_bytes

import (
	"bytes"
	"fmt"
	"testing"
)

// TestBytes_1
// (1)bytes.Buffer与strings.Builder
//	相同点: 都是开箱即用的
//	不同点: strings.Builder只能拼接和导出字符串;而bytes.Buffer不但可以拼接、截断其中的字节序列，以各种形式导出其中的内容，
//     	   还可以顺序地读取其中的子序列。
// (2)bytes.Buffer的注意事项: (重点掌握)
//	a.与strings.Reader类型的Len方法一样，bytes.Buffer的Len方法返回的也是内容容器中未被读取部分的长度，
//    而不是其中已存内容的总长度（以下简称内容长度）
//	b.bytes.Buffer值的容量指的是它的内容容器（也就是那个字节切片）的容量，它只与在当前值之上的写操作有关，并会随着内容的写入而不断增长。
//	c.strings.Reader有一个Size方法可以给出内容长度的值，所以我们用获得的内容长度减去未读部分的长度（Len方法得到），就可以很方便地得到它的已读计数；
//    而bytes.Buffer的Cap方法是内容容器的容量（不是内容长度），故无法计算其已读计数(其结构体封装了已读计数lastRead字段,但是是私有的,没法直接得到)。
func TestBytes_1(t *testing.T) {
	// 示例1。
	var buffer1 bytes.Buffer
	contents := "Simple byte buffer for marshaling data."
	fmt.Printf("Write contents %q ...\n", contents)
	buffer1.WriteString(contents)
	fmt.Printf("The length of buffer: %d\n", buffer1.Len())   // 39
	fmt.Printf("The capacity of buffer: %d\n", buffer1.Cap()) // 64
	fmt.Println()

	// 示例2。
	p1 := make([]byte, 7)
	n, _ := buffer1.Read(p1)
	fmt.Printf("%d bytes were read. (call Read)\n", n)
	fmt.Printf("The length of buffer: %d\n", buffer1.Len())   // 32
	fmt.Printf("The capacity of buffer: %d\n", buffer1.Cap()) // 64
}

// TestBytes_2
// (3)bytes.Buffer类型的值记录的已读计数，在其中起到了怎样的作用?
//	a.读取内容时，相应方法会依据已读计数找到未读部分，并在读取后更新计数。
//  b.写入内容时，如需扩容，相应方法会根据已读计数实现扩容策略。
//	c.截断内容时，相应方法截掉的是已读计数代表索引之后的未读部分。
//	d.读回退时，相应方法需要用已读计数记录回退点。
//	e.重置内容时，相应方法会把已读计数置为0。
//	f.导出内容时，相应方法只会导出已读计数代表的索引之后的未读部分。
//	g.获取长度时，相应方法会依据已读计数和内容容器的长度，计算未读部分的长度并返回。
func TestBytes_2(t *testing.T) {
	// 示例1。
	var contents string
	buffer1 := bytes.NewBufferString(contents)
	fmt.Printf("The length of new buffer with contents %q: %d\n", // 0
		contents, buffer1.Len())
	fmt.Printf("The capacity of new buffer with contents %q: %d\n", // 32
		contents, buffer1.Cap())
	fmt.Println()

	contents = "12345"
	fmt.Printf("Write contents %q ...\n", contents)
	buffer1.WriteString(contents)
	fmt.Printf("The length of buffer: %d\n", buffer1.Len())   // 5
	fmt.Printf("The capacity of buffer: %d\n", buffer1.Cap()) //32
	fmt.Println()

	contents = "67"
	fmt.Printf("Write contents %q ...\n", contents)
	buffer1.WriteString(contents)
	fmt.Printf("The length of buffer: %d\n", buffer1.Len())   // 7
	fmt.Printf("The capacity of buffer: %d\n", buffer1.Cap()) // 32
	fmt.Println()

	contents = "89"
	fmt.Printf("Write contents %q ...\n", contents)
	buffer1.WriteString(contents)
	fmt.Printf("The length of buffer: %d\n", buffer1.Len())   // 9
	fmt.Printf("The capacity of buffer: %d\n", buffer1.Cap()) // 32
	fmt.Print("\n\n")

	// 示例2。
	contents = "abcdefghijk"
	buffer2 := bytes.NewBufferString(contents)
	fmt.Printf("The length of new buffer with contents %q: %d\n", // 11
		contents, buffer2.Len())
	fmt.Printf("The capacity of new buffer with contents %q: %d\n", // 32
		contents, buffer2.Cap())
	fmt.Println()

	n := 10
	fmt.Printf("Grow the buffer with %d ...\n", n)
	buffer2.Grow(n)
	fmt.Printf("The length of buffer: %d\n", buffer2.Len())   // 11
	fmt.Printf("The capacity of buffer: %d\n", buffer2.Cap()) // 32
	fmt.Print("\n\n")

	// 示例3。
	var buffer3 bytes.Buffer
	fmt.Printf("The length of new buffer: %d\n", buffer3.Len())   // 0
	fmt.Printf("The capacity of new buffer: %d\n", buffer3.Cap()) // 0
	fmt.Println()

	contents = "xyz"
	fmt.Printf("Write contents %q ...\n", contents)
	buffer3.WriteString(contents)
	fmt.Printf("The length of buffer: %d\n", buffer3.Len())   // 3
	fmt.Printf("The capacity of buffer: %d\n", buffer3.Cap()) // 64
}

// TestBytes_3
// (4)bytes.Buffer的扩容策略是怎样的?
//	a.Buffer值既可以被手动扩容，也可以进行自动扩容。
//	b.扩容时:
//	  i)如果内容容器的容量与其长度的差，大于或等于另需的字节数，那么扩容代码就会通过切片操作对原有的内容容器的长度进行扩充；
//    ii)如果内容容器的剩余容量不够了，那么扩容代码可能就会用新的内容容器去替代原有的内容容器，从而实现扩容；
//	  iii)如果当前内容容器的容量的一半，仍然大于或等于其现有长度（即未读字节数）再加上另需的字节数的和，那么，扩容代码就会复用现有的内容容器，
//	      并把容器中的未读内容拷贝到它的头部位置（这也意味着其中的已读内容，将会全部被未读内容和之后的新内容覆盖掉）。
//    IV)若不满足iii),扩容代码就只能再创建一个新的内容容器，并把原有容器中的未读内容拷贝进去，最后再用新的容器替换掉原有的容器。
//       这个新容器的容量将会等于原有容量的二倍再加上另需字节数的和。
//	  V)对于处在零值状态的Buffer值来说，如果第一次扩容时的另需字节数不大于64，那么该值就会基于一个预先定义好的、长度为64的字节数组来创建内容容器。
// (5)Q:bytes.Buffer中的哪些方法可能会造成内容的泄露？
//	  A:在bytes.Buffer中，Bytes方法和Next方法都可能会造成内容的泄露。原因在于，它们都把基于内容容器的切片直接返回给了方法的调用方。
//	    不过，如果经过扩容，Buffer值的内容容器或者它的底层数组被重新设定了，那么之前的内容泄露问题就无法再进一步发展了。
//		避免出现内容泄露的最彻底的做法是，在传出切片这类值之前要做好隔离。比如，先对它们进行深度拷贝，然后再把副本传出去。
// (6)Q:对比strings.Builder和bytes.Buffer的String方法，并判断哪一个更高效？原因是什么？
//	  A:strings.Builder的String方法更高效。因为该方法只对其所属值的内容容器（那个字节切片）做了简单的类型转换，并且直接使用了底层的值（或者说内存空间）。
//	    它的源码如下：
//		// String returns the accumulated string.
//		func (b *Builder) String() string {
//			return *(*string)(unsafe.Pointer(&b.buf))
//		}
//		数组值和字符串值在底层的存储方式其实是一样的。所以从切片值到字符串值的指针值的转换可以是直截了当的。又由于字符串值是不可变的，所以这样做也是安全的。
//		不过，由于一些历史、结构和功能方面的原因，bytes.Buffer的String方法却不能这样做。
func TestBytes_3(t *testing.T) {
	// 示例1。
	contents := "ab"
	buffer1 := bytes.NewBufferString(contents)
	fmt.Printf("The capacity of new buffer with contents %q: %d\n", // 8
		contents, buffer1.Cap())
	fmt.Println()

	unreadBytes := buffer1.Bytes()
	fmt.Printf("The unread bytes of the buffer: %v\n", unreadBytes) // [97 98]
	fmt.Println()

	contents = "cdefg"
	fmt.Printf("Write contents %q ...\n", contents)
	buffer1.WriteString(contents)
	fmt.Printf("The capacity of buffer: %d\n", buffer1.Cap()) // 8
	fmt.Println()

	// 只要扩充一下之前拿到的未读字节切片unreadBytes，
	// 就可以用它来读取甚至修改buffer中的后续内容。
	unreadBytes = unreadBytes[:cap(unreadBytes)]
	fmt.Printf("The unread bytes of the buffer: %v\n", unreadBytes) // [97 98 99 100 101 102 103 0]
	fmt.Println()

	value := byte('X')
	fmt.Printf("Set a byte in the unread bytes to %v ...\n", value)
	unreadBytes[len(unreadBytes)-2] = value
	fmt.Printf("The unread bytes of the buffer: %v\n", buffer1.Bytes()) // [97 98 99 100 101 102 88]
	fmt.Println()

	// 不过，在buffer的内容容器真正扩容之后就无法这么做了。
	contents = "hijklmn"
	fmt.Printf("Write contents %q ...\n", contents)
	buffer1.WriteString(contents)
	fmt.Printf("The capacity of buffer: %d\n", buffer1.Cap()) // 23
	fmt.Println()

	unreadBytes = unreadBytes[:cap(unreadBytes)]
	fmt.Printf("The unread bytes of the buffer: %v\n", unreadBytes) //  [97 98 99 100 101 102 88 0]
	fmt.Print("\n\n")

	// 示例2。
	// Next方法返回的后续字节切片也存在相同的问题。
	contents = "12"
	buffer2 := bytes.NewBufferString(contents)
	fmt.Printf("The capacity of new buffer with contents %q: %d\n", // 8
		contents, buffer2.Cap())
	fmt.Println()

	nextBytes := buffer2.Next(2)
	fmt.Printf("The next bytes of the buffer: %v\n", nextBytes) // [49 50]
	fmt.Println()

	contents = "34567"
	fmt.Printf("Write contents %q ...\n", contents)
	buffer2.WriteString(contents)
	fmt.Printf("The capacity of buffer: %d\n", buffer2.Cap()) // 8
	fmt.Println()

	// 只要扩充一下之前拿到的后续字节切片nextBytes，
	// 就可以用它来读取甚至修改buffer中的后续内容。
	nextBytes = nextBytes[:cap(nextBytes)]
	fmt.Printf("The next bytes of the buffer: %v\n", nextBytes) // [49 50 51 52 53 54 55 0]
	fmt.Println()

	value = byte('X')
	fmt.Printf("Set a byte in the next bytes to %v ...\n", value)
	nextBytes[len(nextBytes)-2] = value
	fmt.Printf("The unread bytes of the buffer: %v\n", buffer2.Bytes()) // [51 52 53 54 88]
	fmt.Println()

	// 不过，在buffer的内容容器真正扩容之后就无法这么做了。
	contents = "89101112"
	fmt.Printf("Write contents %q ...\n", contents)
	buffer2.WriteString(contents)
	fmt.Printf("The capacity of buffer: %d\n", buffer2.Cap()) // 24
	fmt.Println()

	nextBytes = nextBytes[:cap(nextBytes)]
	fmt.Printf("The next bytes of the buffer: %v\n", nextBytes) // [49 50 51 52 53 54 88 0]
}
