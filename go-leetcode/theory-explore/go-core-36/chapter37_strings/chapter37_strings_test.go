package chapter37_strings

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

// TestStrings_1
// (1)Q:与string值相比，strings.Builder类型的值有哪些优势?
//	  A:strings.Builder类型的值（以下简称Builder值）的优势有下面的三种：（与string值相比，Builder值的优势其实主要体现在字符串拼接方面）
//		a.已存在的内容不可变，但可以拼接更多的内容。
//		b.减少了内存分配和内容拷贝的次数。
//		c.可将内容重置，可重用值
// (2)关于string,有如下特点:
//	a.string类型的值是不可变的。
//    如果我们想获得一个不一样的字符串，那么就只能基于原字符串进行裁剪(使用切片表达式)、拼接(使用操作符+)等操作，从而生成一个新的字符串。
//	b.值会存储到一块连续的内存空间。
//	c.字节数量即为string值的长度。
//	d.在一个string值上应用切片表达式，就相当于在对其底层的字节数组做切片。
//	e.在进行字符串拼接的时候，Go语言会把所有被拼接的字符串依次拷贝到一个崭新且足够大的连续内存空间中，并把持有相应指针值的string值作为结果返回。
// (3)关于Builder,有如下特点:
//	a.Builder值中的内容是 不可变 的;Builder底层切片中的内容只能够被拼接或者完全重置。
//	b.Builder值会自动地对自身的内容容器 进行扩容 ;除了Builder值的自动扩容，我们还可以选择手动扩容(Grow方法)。
//	  i)在向Builder值拼接内容的时候并不一定会引起扩容。只要内容容器的容量够用，扩容就不会进行，针对于此的内存分配也不会发生。
//	  同时，只要没有扩容，Builder值中已存在的内容就不会再被拷贝。
//	  ii)Grow方法生成一个字节切片作为新的内容容器，该切片的容量会是原容器容量的二倍再加上n。
//	  iii)Grow方法还可能什么都不做。这种情况的前提条件是：
//	 	  当前的内容容器中的未用容量已经够用了，即：未用容量大于或等于n。这里的前提条件与前面提到的自动扩容策略中的前提条件是类似的。
// c.Builder值是可以被重用的。通过调用它的Reset方法,我们可以让Builder值重新回到零值状态，就像它从未被使用过那样。
//   原有的内容容器会被直接丢弃，将会被垃圾回收器标记并回收掉。
func TestStrings_1(t *testing.T) {
	// 示例1。
	var builder1 strings.Builder
	builder1.WriteString("A Builder is used to efficiently build a string using Write methods.")
	fmt.Printf("The first output(%d):\n%q\n", builder1.Len(), builder1.String()) // 68
	fmt.Println()
	builder1.WriteByte(' ')
	builder1.WriteString("It minimizes memory copying. The zero value is ready to use.")
	builder1.Write([]byte{'\n', '\n'})
	builder1.WriteString("Do not copy a non-zero Builder.")
	fmt.Printf("The second output(%d):\n\"%s\"\n", builder1.Len(), builder1.String()) // 162
	fmt.Println()

	// iii)Grow方法还可能什么都不做。这种情况的前提条件是：
	//	   当前的内容容器中的未用容量已经够用了，即：未用容量大于或等于n。这里的前提条件与前面提到的自动扩容策略中的前提条件是类似的。
	// 示例2。
	fmt.Println("Grow the builder ...")
	builder1.Grow(10)
	fmt.Printf("The length of contents in the builder is %d.\n", builder1.Len()) // 162
	fmt.Println()

	// c.Builder值是可以被重用的。通过调用它的Reset方法,我们可以让Builder值重新回到零值状态，就像它从未被使用过那样。
	//   原有的内容容器会被直接丢弃，将会被垃圾回收器标记并回收掉。
	// 示例3。
	fmt.Println("Reset the builder ...")
	builder1.Reset()
	fmt.Printf("The third output(%d):\n%q\n", builder1.Len(), builder1.String()) // 0
}

// TestStrings_2
// (4)Q:strings.Builder类型在使用上有约束吗?
//	  A:有约束，概括如下：
//		a.在已被真正使用后就不可再被复制(这里所说的复制方式，包括但不限于在函数间传递值、通过通道传递值、把值赋予变量等；否则就会出现panic)；但是可以传递指针。
//		  复制不会panic的两种特例:
//			i)我们还可以先使用再传递(即复制)，只要在传递之前调用它的Reset方法即可。
//			ii)对于处在零值状态的Builder值，复制不会有任何问题。
//		b.由于其内容不是完全不可变的，所以需要使用方自行解决操作冲突和并发安全问题(传指针的时候就要考虑并发安全问题)。
//		  我们在通过传递其指针值共享Builder值的时候，一定要确保各方对它的使用是正确、有序的，并且是并发安全的；
//		  而最彻底的解决方案是，绝不共享Builder值以及它的指针值。

func TestStrings_2(t *testing.T) {
	// a.在已被真正使用后就不可再被复制(这里所说的复制方式，包括但不限于在函数间传递值、通过通道传递值、把值赋予变量等；否则就会出现panic)；但是可以传递指针。
	// 示例1。
	var builder1 strings.Builder
	builder1.Grow(1)

	f1 := func(b strings.Builder) {
		//b.Grow(1) // 这里会引发panic。
	}
	f1(builder1)

	ch1 := make(chan strings.Builder, 1)
	ch1 <- builder1
	builder2 := <-ch1
	//builder2.Grow(1) // 这里会引发panic。
	_ = builder2

	builder3 := builder1
	//builder3.Grow(1) // 这里会引发panic。
	_ = builder3

	// 示例2。
	f2 := func(bp *strings.Builder) {
		(*bp).Grow(1) // 这里虽然不会引发panic，但不是并发安全的。
		builder4 := *bp
		//builder4.Grow(1) // 这里会引发panic。
		_ = builder4
	}
	f2(&builder1)

	// i)我们还可以先使用再传递(即复制)，只要在传递之前调用它的Reset方法即可。
	builder1.Reset()
	builder5 := builder1
	builder5.Grow(1) // 这里不会引发panic。
}

// TestStrings_3
// (5)Q:为什么说strings.Reader类型的值可以高效地读取字符串?
//    A:主要原因是在读取的过程中，Reader值会保存已读取的字节的计数。
//		具体而言:
//		a.已读计数也代表着下一次读取的起始索引位置（依靠已读计数以及针对字符串值的切片表达式，从而实现快速读取）。
//		b.已读计数也是读取回退和位置设定时的重要依据。已读计数=reader.Size()-reader.Len()
//		  其中reader.Size()代表容量;reader.Len()代表未读的数量
//		c.Reader值拥有的大部分用于读取的方法(如ReadByte方法,ReadRune方法)都会及时地更新已读计数；
//        ReadAt方法算是一个例外。它既不会依据已读计数进行读取，也不会在读取后更新它；所以该方法可以自由地读取其所属的Reader值中的任何内容。
//		  Reader值的Seek方法也会更新该值的已读计数；实际上，这个Seek方法的主要作用正是设定下一次读取的起始索引位置。
//		  如果我们把常量io.SeekCurrent的值作为第二个参数值传给该方法，那么它还会依据当前的已读计数，以及第一个参数offset的值来计算新的计数值。
// (6)Q:*strings.Builder和*strings.Reader都分别实现了哪些接口？这样做有什么好处吗？
//	  A:strings.Builder类型实现了 3 个接口，分别是：fmt.Stringer、io.Writer和io.ByteWriter。
//		strings.Reader类型则实现了 8 个接口，即：io.Reader、io.ReaderAt、io.ByteReader、io.RuneReader、io.Seeker、io.ByteScanner、io.RuneScanner和io.WriterTo。
//		好处是显而易见的。实现的接口越多，它们的用途就越广。它们会适用于那些要求参数的类型为这些接口类型的地方。
func TestStrings_3(t *testing.T) {
	//a.已读计数也代表着下一次读取的起始索引位置（依靠已读计数以及针对字符串值的切片表达式，从而实现快速读取）。
	//b.已读计数也是读取回退和位置设定时的重要依据。已读计数=reader.Size()-reader.Len()
	//  其中reader.Size()代表容量;reader.Len()代表未读的数量
	// 示例1。
	reader1 := strings.NewReader(
		"NewReader returns a new Reader reading from s. " +
			"It is similar to bytes.NewBufferString but more efficient and read-only.")
	fmt.Printf("The size of reader: %d\n", reader1.Size()) // 119
	fmt.Printf("The reading index in reader: %d\n",        // 0
		reader1.Size()-int64(reader1.Len()))

	buf1 := make([]byte, 47)
	n, _ := reader1.Read(buf1)
	fmt.Printf("%d bytes were read. (call Read)\n", n)
	fmt.Printf("The reading index in reader: %d\n", // 47
		reader1.Size()-int64(reader1.Len()))
	fmt.Println()

	// ReadAt方法算是一个例外。它既不会依据已读计数进行读取，也不会在读取后更新它；所以该方法可以自由地读取其所属的Reader值中的任何内容。
	// 示例2。
	buf2 := make([]byte, 21)
	offset1 := int64(64)
	n, _ = reader1.ReadAt(buf2, offset1)
	fmt.Printf("%d bytes were read. (call ReadAt, offset: %d)\n", n, offset1)
	fmt.Printf("The reading index in reader: %d\n", // 47
		reader1.Size()-int64(reader1.Len()))
	fmt.Println()

	//Reader值的Seek方法也会更新该值的已读计数；实际上，这个Seek方法的主要作用正是设定下一次读取的起始索引位置。
	//如果我们把常量io.SeekCurrent的值作为第二个参数值传给该方法，那么它还会依据当前的已读计数，以及第一个参数offset的值来计算新的计数值。
	// 示例3。
	offset2 := int64(17)
	expectedIndex := reader1.Size() - int64(reader1.Len()) + offset2
	fmt.Printf("Seek with offset %d and whence %d ...\n", offset2, io.SeekCurrent)
	readingIndex, _ := reader1.Seek(offset2, io.SeekCurrent)
	fmt.Printf("The reading index in reader: %d (returned by Seek)\n", readingIndex) // 47+17=64
	fmt.Printf("The reading index in reader: %d (computed by me)\n", expectedIndex)  // 64

	n, _ = reader1.Read(buf2)
	fmt.Printf("%d bytes were read. (call Read)\n", n)
	fmt.Printf("The reading index in reader: %d\n",
		reader1.Size()-int64(reader1.Len())) // 64+21=85
}
