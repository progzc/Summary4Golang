package chapter40_41_io

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestIo_1
// (1)Q:strings.Builder、strings.Reader和bytes.Buffer都实现了哪些接口?
//
//		  A:
//		   a.strings.Builder类型主要用于构建字符串,它的指针类型实现的接口有io.Writer、io.ByteWriter和fmt.Stringer。
//		     另外，它其实还实现了一个io包的包级私有接口io.stringWriter（自 Go 1.12 起它会更名为io.StringWriter）。
//		   b.strings.Reader类型主要用于读取字符串,它的指针类型实现的接口比较多,包括以下8个接口:
//			 	io.Reader;		io.ReaderAt;		io.ByteReader;		io.RuneReader;
//			 	io.Seeker;		io.ByteScanner;		io.RuneScanner;		io.WriterTo;
//			 其中，io.ByteScanner是io.ByteReader的扩展接口，而io.RuneScanner又是io.RuneReader的扩展接口。
//		   c.bytes.Buffer是集读、写功能于一身的数据类型,它非常适合作为字节序列的缓冲区。
//		     该指针类型实现的读取相关的接口有下面6个:
//			 	io.Reader;			io.ByteReader;		io.RuneReader;		io.ByteScanner;
//			 	io.RuneScanner;		io.WriterTo
//	      该指针类型实现的写入相关的接口则有以下4个:
//	      	io.Writer;			io.ByteWriter;		io.stringWriter;	io.ReaderFrom;
//			 该指针类型实现的导出相关的接口有下面1个:
//			 	fmt.Stringer
//
// (2)Q:io包中接口的好处与优势?
//
//	  A:为了提高不同程序实体之间的互操作性。换句话说，如此一来，Go 语言的各种库中，能够操作它们的函数和数据类型明显多了很多。
//		此外,在Go语言中，对接口的扩展是通过接口类型之间的嵌入来实现的，这也常被叫做接口的组合。Go 语言提倡使用小接口加接口组合的方式，
//		来扩展程序的行为以及增加程序的灵活性。io代码包恰恰就可以作为这样的一个标杆，它可以成为我们运用这种技巧时的一个参考标准。
func TestIo_1(t *testing.T) {
	// 示例1。
	builder := new(strings.Builder)
	_ = interface{}(builder).(io.Writer)
	_ = interface{}(builder).(io.ByteWriter)
	_ = interface{}(builder).(fmt.Stringer)

	// 示例2。
	reader := strings.NewReader("")
	_ = interface{}(reader).(io.Reader)
	_ = interface{}(reader).(io.ReaderAt)
	_ = interface{}(reader).(io.ByteReader)
	_ = interface{}(reader).(io.RuneReader)
	_ = interface{}(reader).(io.Seeker)
	_ = interface{}(reader).(io.ByteScanner)
	_ = interface{}(reader).(io.RuneScanner)
	_ = interface{}(reader).(io.WriterTo)

	// 示例3。
	buffer := bytes.NewBuffer([]byte{})
	_ = interface{}(buffer).(io.Reader)
	_ = interface{}(buffer).(io.ByteReader)
	_ = interface{}(buffer).(io.RuneReader)
	_ = interface{}(buffer).(io.ByteScanner)
	_ = interface{}(buffer).(io.RuneScanner)
	_ = interface{}(buffer).(io.WriterTo)

	_ = interface{}(buffer).(io.Writer)
	_ = interface{}(buffer).(io.ByteWriter)
	_ = interface{}(buffer).(io.ReaderFrom)

	_ = interface{}(buffer).(fmt.Stringer)

	// 示例4。
	src := strings.NewReader(
		"CopyN copies n bytes (or until an error) from src to dst. " +
			"It returns the number of bytes copied and " +
			"the earliest error encountered while copying.")
	dst := new(strings.Builder)
	written, err := io.CopyN(dst, src, 58)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	} else {
		fmt.Printf("Written(%d): %q\n", written, dst.String())
	}
}

// TestIo_2
// (3)Q:在io包中，io.Reader的扩展接口和实现类型都有哪些？它们分别都有什么功用？
//
//	  A:在io包中，io.Reader的扩展接口有下面几种:
//		a.io.ReadWriter：此接口既是io.Reader的扩展接口，也是io.Writer的扩展接口。
//		b.io.ReadCloser：此接口除了包含基本的字节序列读取方法之外，还拥有一个基本的关闭方法Close。
//		c.io.ReadWriteCloser：很明显，此接口是io.Reader、io.Writer和io.Closer这三个接口的组合。
//		d.io.ReadSeeker：此接口的特点是拥有一个用于寻找读写位置的基本方法Seek。更具体地说，
//	      该方法可以根据给定的偏移量基于数据的起始位置、末尾位置，或者当前读写位置去寻找新的读写位置。
//		e.io.ReadWriteSeeker：显然，此接口是另一个三合一的扩展接口，它是io.Reader、io.Writer和io.Seeker的组合。
//		io包中的io.Reader接口的实现类型:
//		a.*io.LimitedReader：此类型的基本类型会包装io.Reader类型的值，并提供一个额外的受限读取的功能。
//		b.*io.SectionReader：此类型的基本类型可以包装io.ReaderAt类型的值，并且会限制它的Read方法，只能够读取原始数据中的某一个部分（或者说某一段）。
//		c.*io.teeReader：此类型是一个包级私有的数据类型，也是io.TeeReader函数结果值的实际类型。
//		d.*io.multiReader：此类型也是一个包级私有的数据类型。
//		e.*io.pipe：此类型为一个包级私有的数据类型，它比上述类型都要复杂得多。
//		f.*io.PipeReader：此类型可以被视为io.pipe类型的代理类型。
func TestIo_2(t *testing.T) {
	comment := "Package io provides basic interfaces to I/O primitives. " +
		"Its primary job is to wrap existing implementations of such primitives, " +
		"such as those in package os, " +
		"into shared public interfaces that abstract the functionality, " +
		"plus some other related primitives."

	// 示例1。
	fmt.Println("New a string reader and name it \"reader1\" ...")
	reader1 := strings.NewReader(comment)
	buf1 := make([]byte, 7)
	n, err := reader1.Read(buf1)
	var offset1, index1 int64
	executeIfNoErr(err, func() {
		fmt.Printf("Read(%d): %q\n", n, buf1[:n])
		offset1 = int64(53)
		index1, err = reader1.Seek(offset1, io.SeekCurrent)
	})
	executeIfNoErr(err, func() {
		fmt.Printf("The new index after seeking from current with offset %d: %d\n",
			offset1, index1)
		n, err = reader1.Read(buf1)
	})
	executeIfNoErr(err, func() {
		fmt.Printf("Read(%d): %q\n", n, buf1[:n])
	})
	fmt.Println()

	// 示例2。
	reader1.Reset(comment)
	num1 := int64(7)
	fmt.Printf("New a limited reader with reader1 and number %d ...\n", num1)
	reader2 := io.LimitReader(reader1, 7)
	buf2 := make([]byte, 10)
	for i := 0; i < 3; i++ {
		n, err = reader2.Read(buf2)
		executeIfNoErr(err, func() {
			fmt.Printf("Read(%d): %q\n", n, buf2[:n])
		})
	}
	fmt.Println()

	// 示例3。
	reader1.Reset(comment)
	offset2 := int64(56)
	num2 := int64(72)
	fmt.Printf("New a section reader with reader1, offset %d and number %d ...\n", offset2, num2)
	reader3 := io.NewSectionReader(reader1, offset2, num2)
	buf3 := make([]byte, 20)
	for i := 0; i < 5; i++ {
		n, err = reader3.Read(buf3)
		executeIfNoErr(err, func() {
			fmt.Printf("Read(%d): %q\n", n, buf3[:n])
		})
	}
	fmt.Println()

	// 示例4。
	reader1.Reset(comment)
	writer1 := new(strings.Builder)
	fmt.Println("New a tee reader with reader1 and writer1 ...")
	reader4 := io.TeeReader(reader1, writer1)
	buf4 := make([]byte, 40)
	for i := 0; i < 8; i++ {
		n, err = reader4.Read(buf4)
		executeIfNoErr(err, func() {
			fmt.Printf("Read(%d): %q\n", n, buf4[:n])
		})
	}
	fmt.Println()

	// 示例5。
	reader5a := strings.NewReader(
		"MultiReader returns a Reader that's the logical concatenation of " +
			"the provided input readers.")
	reader5b := strings.NewReader("They're read sequentially.")
	reader5c := strings.NewReader("Once all inputs have returned EOF, " +
		"Read will return EOF.")
	reader5d := strings.NewReader("If any of the readers return a non-nil, " +
		"non-EOF error, Read will return that error.")
	fmt.Println("New a multi-reader with 4 readers ...")
	reader5 := io.MultiReader(reader5a, reader5b, reader5c, reader5d)
	buf5 := make([]byte, 50)
	for i := 0; i < 8; i++ {
		n, err = reader5.Read(buf5)
		executeIfNoErr(err, func() {
			fmt.Printf("Read(%d): %q\n", n, buf5[:n])
		})
	}
	fmt.Println()

	// 示例6。
	fmt.Println("New a synchronous in-memory pipe ...")
	pReader, pWriter := io.Pipe()
	_ = interface{}(pReader).(io.ReadCloser)
	_ = interface{}(pWriter).(io.WriteCloser)

	comments := [][]byte{
		[]byte("Pipe creates a synchronous in-memory pipe."),
		[]byte("It can be used to connect code expecting an io.Reader "),
		[]byte("with code expecting an io.Writer."),
	}

	// 这里添加这个同步工具纯属为了保证下面示例中的打印语句都能够执行完成。
	// 在实际使用中没有必要这样做。
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for _, d := range comments {
			time.Sleep(time.Millisecond * 500)
			n, err := pWriter.Write(d)
			if err != nil {
				fmt.Printf("write error: %v\n", err)
				break
			}
			fmt.Printf("Written(%d): %q\n", n, d)
		}
		pWriter.Close()
	}()
	go func() {
		defer wg.Done()
		wBuf := make([]byte, 55)
		for {
			n, err := pReader.Read(wBuf)
			if err != nil {
				fmt.Printf("read error: %v\n", err)
				break
			}
			fmt.Printf("Read(%d): %q\n", n, wBuf[:n])
		}
	}()
	wg.Wait()
}

func executeIfNoErr(err error, f func()) {
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	f()
}

// TestIo_3
// (4)Q:io包中的接口都有哪些？它们之间都有着怎样的关系？
//
//	  A:简单接口:没有嵌入其他接口并且只定义了一个方法的接口。在io包中，这样的接口一共有 11 个。
//				可以把io包中的简单接口分为四大类。这四大类接口分别针对于四种操作，即：读取、写入、关闭和读写位置设定。
//				前三种操作属于基本的 I/O 操作。
//		核心接口:有着众多的扩展接口和实现类型的接口。
//	            io包中的核心接口只有 3 个，它们是：io.Reader、io.Writer和io.Closer。
//
// (5)Q:io包中的同步内存管道的运作机制是什么？
//
//	  A:io.Pipe函数会返回一个io.PipeReader类型的值和一个io.PipeWriter类型的值，并将它们分别作为管道的两端；
//	    而这两个值在底层其实只是代理了同一个*io.pipe类型值的功能而已。
//		io.pipe类型通过无缓冲的通道实现了读操作与写操作之间的同步，并且通过互斥锁实现了写操作之间的串行化。另外，它还使用原子值来处理错误；
//		这些共同保证了这个同步内存管道的并发安全性。
func TestIo_3(t *testing.T) {
	comment := "Because these interfaces and primitives wrap lower-level operations with various implementations, " +
		"unless otherwise informed clients should not assume they are safe for parallel execution."
	basicReader := strings.NewReader(comment)
	basicWriter := new(strings.Builder)

	// 示例1。
	reader1 := io.LimitReader(basicReader, 98)
	_ = interface{}(reader1).(io.Reader)

	// 示例2。
	reader2 := io.NewSectionReader(basicReader, 98, 89)
	_ = interface{}(reader2).(io.Reader)
	_ = interface{}(reader2).(io.ReaderAt)
	_ = interface{}(reader2).(io.Seeker)

	// 示例3。
	reader3 := io.TeeReader(basicReader, basicWriter)
	_ = interface{}(reader3).(io.Reader)

	// 示例4。
	reader4 := io.MultiReader(reader1)
	_ = interface{}(reader4).(io.Reader)

	// 示例5。
	writer1 := io.MultiWriter(basicWriter)
	_ = interface{}(writer1).(io.Writer)

	// 示例6。
	pReader, pWriter := io.Pipe()
	_ = interface{}(pReader).(io.Reader)
	_ = interface{}(pReader).(io.Closer)
	_ = interface{}(pWriter).(io.Writer)
	_ = interface{}(pWriter).(io.Closer)
}
