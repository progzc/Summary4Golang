package chapter42_43_bufio

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// TestBufIo_1
// (1)bufio包中的数据类型
//	  Reader;	Scanner;	Writer;		ReadWriter;
// (2)Q:bufio.Reader类型值中的缓冲区起着怎样的作用?
//	  A:bufio.Reader类型的值（以下简称Reader值）内的缓冲区，其实就是一个数据存储中介，它介于底层读取器与读取方法及其调用方之间。
//	    所谓的底层读取器，就是在初始化此类值的时候传入的io.Reader类型的参数值。
//		Reader值的读取方法一般都会先从其所属值的缓冲区中读取数据。同时，在必要的时候，它们还会预先从底层读取器那里读出一部分数据，并暂存于缓冲区之中以备后用。
// (2)关于bufio.Reader
//	a.bufio.Reader类型并不是开箱即用
//	b.包含的字段:
//		buf：[]byte类型的字段，即字节切片，代表缓冲区。虽然它是切片类型的，但是其长度却会在初始化的时候指定，并在之后保持不变。
//		rd：io.Reader类型的字段，代表底层读取器。缓冲区中的数据就是从这里拷贝来的。
//		r：int类型的字段，代表对缓冲区进行下一次读取时的开始索引。我们可以称它为已读计数。
//		w：int类型的字段，代表对缓冲区进行下一次写入时的开始索引。我们可以称之为已写计数。
//		err：error类型的字段。它的值用于表示在从底层读取器获得数据时发生的错误。这里的值在被读取或忽略之后，该字段会被置为nil。
//		lastByte：int类型的字段，用于记录缓冲区中最后一个被读取的字节。读回退时会用到它的值。
//		astRuneSize：int类型的字段，用于记录缓冲区中最后一个被读取的 Unicode 字符所占用的字节数。读回退的时候会用到它的值。
//	c.在bufio.Reader类型拥有的读取方法中，Peek方法和ReadSlice方法都会调用该类型一个名为fill的包级私有方法。
//    fill方法的作用是填充内部缓冲区(fill会涉及到缓冲区的压缩)，对缓冲区的压缩包括两个步骤：
//   	第一步，把缓冲区中在[已读计数, 已写计数)范围之内的所有元素值（或者说字节）都依次拷贝到缓冲区的头部。
//		第二步，fill方法会把已写计数的新值设定为原已写计数与原已读计数的差。这个差所代表的索引，就是压缩后第一次写入字节时的开始索引。
//	d.fill方法只要在开始时发现其所属值的已读计数大于0，就会对缓冲区进行一次压缩。之后，如果缓冲区中还有可写的位置，那么该方法就会对其进行填充。
//	  在填充缓冲区的时候，fill方法会试图从底层读取器那里，读取足够多的字节，并尽量把从已写计数代表的索引位置到缓冲区末尾之间的空间都填满。
//	  fill方法会及时地更新已写计数，以保证填充的正确性和顺序性。
// (3)关于bufio.Writer
//	a.包含的字段
//		err：error类型的字段。它的值用于表示在向底层写入器写数据时发生的错误。
//		buf：[]byte类型的字段，代表缓冲区。在初始化之后，它的长度会保持不变。
//		n：int类型的字段，代表对缓冲区进行下一次写入时的开始索引。我们可以称之为已写计数。
//		wr：io.Writer类型的字段，代表底层写入容器。
//	b.bufio.Writer类型值中缓冲的数据什么时候会被写到它的底层写入器?
//		Flush方法：主要功能是把相应缓冲区中暂存的所有数据，都写到底层写入器中。数据一旦被写进底层写入器，该方法就会把它们从缓冲区中删除（指逻辑上的删除）掉。
//				  bufio.Writer类型值（以下简称Writer值）拥有的所有数据写入方法（如Write、WriteByte、WriteRune方法）都会在 必要的时候 调用它的Flush方法。
//				  不过，在你把所有的数据都写入Writer值之后，再调用一下它的Flush方法，显然是最稳妥的。
func TestBufIo_1(t *testing.T) {
	comment := "Package bufio implements buffered I/O. " +
		"It wraps an io.Reader or io.Writer object, " +
		"creating another object (Reader or Writer) that " +
		"also implements the interface but provides buffering and " +
		"some help for textual I/O."
	basicReader := strings.NewReader(comment)
	fmt.Printf("The size of basic reader: %d\n", basicReader.Size()) // 213
	fmt.Println()

	// 示例1。
	fmt.Println("New a buffered reader ...")
	reader1 := bufio.NewReader(basicReader)
	fmt.Printf("The default size of buffered reader: %d\n", reader1.Size()) // 4096
	// 此时reader1的缓冲区还没有被填充。
	fmt.Printf("The number of unread bytes in the buffer: %d\n", reader1.Buffered()) // 0
	fmt.Println()

	// 示例2。
	bytes, err := reader1.Peek(7)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Printf("Peeked contents(%d): %q\n", len(bytes), bytes)
	fmt.Printf("The number of unread bytes in the buffer: %d\n", reader1.Buffered()) // 213
	fmt.Println()

	// 示例3。
	buf1 := make([]byte, 7)
	n, err := reader1.Read(buf1)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Printf("Read contents(%d): %q\n", n, buf1)
	fmt.Printf("The number of unread bytes in the buffer: %d\n", reader1.Buffered()) // 206
	fmt.Println()

	// 示例4。
	fmt.Printf("Reset the basic reader (size: %d) ...\n", len(comment))
	basicReader.Reset(comment)
	fmt.Printf("Reset the buffered reader (size: %d) ...\n", reader1.Size()) // 4096
	reader1.Reset(basicReader)
	peekNum := len(comment) + 1
	fmt.Printf("Peek %d bytes ...\n", peekNum) // 214
	bytes, err = reader1.Peek(peekNum)
	if err != nil {
		fmt.Printf("error: %v\n", err) // EOF
	}
	fmt.Printf("The number of peeked bytes: %d\n", len(bytes)) // 213
	fmt.Println()

	// 示例5。
	fmt.Printf("Reset the basic reader (size: %d) ...\n", len(comment))
	basicReader.Reset(comment)
	size := 300
	fmt.Printf("New a buffered reader with size %d ...\n", size)
	reader2 := bufio.NewReaderSize(basicReader, size)
	peekNum = size + 1
	fmt.Printf("Peek %d bytes ...\n", peekNum)
	bytes, err = reader2.Peek(peekNum)
	if err != nil {
		fmt.Printf("error: %v\n", err) // bufio: buffer full
	}
	fmt.Printf("The number of peeked bytes: %d\n", len(bytes)) // 213
}

// TestBufIo_2
// (4)Q:bufio.Reader类型读取方法有哪些不同？主要的读取方法是Peek、Read、ReadSlice和ReadBytes
//	a.Peek方法：读取并返回其缓冲区中的n个未读字节，并且它会从已读计数代表的索引位置开始读；Peek 方法有一个鲜明的特点，那就是：
//	           即使它读取了缓冲区中的数据，也不会更改已读计数的值。
//	b.Read方法
//	c.ReadSlice方法：持续地读取数据，直至遇到调用方给定的分隔符为止。
//	d.ReadBytes方法：持续地读取数据，直至遇到调用方给定的分隔符为止。
//	e.ReadLine方法：
// (5)关于内容泄露
//	a.bufio.Reader类型的Peek方法、ReadSlice方法和ReadLine方法都有可能会造成内容泄露,因为它们在正常的情况下都会返回直接基于缓冲区的字节切片。
// (6)Q:bufio.Scanner类型的主要功用是什么？它有哪些特点？
//	  A:bufio.Scanner类型俗称带缓存的扫描器。
//		a.我们可以自定义每次扫描的边界，或者说内容的分段方法。在默认情况下，扫描器会以行为单位对目标内容进行扫描。bufio代码包提供了一些现成的分段方法。
//	      实际上，扫描器在默认情况下会使用bufio.ScanLines函数作为分段方法。
//		b。我们还可以在扫描之前自定义缓存的载体和缓存的最大容量，这需要调用它的Buffer方法。在默认情况下，扫描器内部设定的最大缓存容量是64K个字节。
func TestBufIo_2(t *testing.T) {
	comment := "Package bufio implements buffered I/O. " +
		"It wraps an io.Reader or io.Writer object, " +
		"creating another object (Reader or Writer) that " +
		"also implements the interface but provides buffering and " +
		"some help for textual I/O."
	basicReader := strings.NewReader(comment)
	fmt.Printf("The size of basic reader: %d\n", basicReader.Size()) // 213

	size := 300
	fmt.Printf("New a buffered reader with size %d ...\n", size)
	reader1 := bufio.NewReaderSize(basicReader, size)
	fmt.Println()

	fmt.Print("[ About 'Peek' method ]\n\n")
	// 示例1。
	peekNum := 38
	fmt.Printf("Peek %d bytes ...\n", peekNum)
	bytes, err := reader1.Peek(peekNum)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Printf("Peeked contents(%d): %q\n", len(bytes), bytes)
	fmt.Printf("The number of unread bytes in the buffer: %d\n", reader1.Buffered()) // 213
	fmt.Println()

	fmt.Print("[ About 'Read' method ]\n\n")
	// 示例2。
	readNum := 38
	buf1 := make([]byte, readNum)
	fmt.Printf("Read %d bytes ...\n", readNum)
	n, err := reader1.Read(buf1)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Printf("Read contents(%d): %q\n", n, buf1)
	fmt.Printf("The number of unread bytes in the buffer: %d\n", reader1.Buffered()) // 175
	fmt.Println()

	fmt.Print("[ About 'ReadSlice' method ]\n\n")
	// 示例3。
	fmt.Println("Reset the basic reader ...")
	basicReader.Reset(comment)
	fmt.Println("Reset the buffered reader ...")
	reader1.Reset(basicReader)
	fmt.Println()

	delimiter := byte('(')
	fmt.Printf("Read slice with delimiter %q...\n", delimiter)
	line, err := reader1.ReadSlice(delimiter)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Printf("Read contents(%d): %q\n", len(line), line)
	fmt.Printf("The number of unread bytes in the buffer: %d\n", reader1.Buffered())
	fmt.Println()

	delimiter = byte('[')
	fmt.Printf("Read slice with delimiter %q...\n", delimiter)
	line, err = reader1.ReadSlice(delimiter)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Printf("Read contents(%d): %q\n", len(line), line)
	fmt.Printf("The number of unread bytes in the buffer: %d\n", reader1.Buffered())
	fmt.Println()

	// 示例4。
	fmt.Println("Reset the basic reader ...")
	basicReader.Reset(comment)
	size = 200
	fmt.Printf("New a buffered reader with size %d ...\n", size)
	reader2 := bufio.NewReaderSize(basicReader, size)
	fmt.Println()

	delimiter = byte('[')
	fmt.Printf("Read slice with delimiter %q...\n", delimiter)
	line, err = reader2.ReadSlice(delimiter)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Printf("Read contents(%d): %q\n", len(line), line)
	fmt.Printf("The number of unread bytes in the buffer: %d\n", reader2.Buffered())
	fmt.Println()

	fmt.Print("[ About 'ReadBytes' method ]\n\n")
	// 示例5。
	fmt.Println("Reset the basic reader ...")
	basicReader.Reset(comment)
	size = 200
	fmt.Printf("New a buffered reader with size %d ...\n", size)
	reader3 := bufio.NewReaderSize(basicReader, size)
	fmt.Println()

	delimiter = byte('[')
	fmt.Printf("Read bytes with delimiter %q...\n", delimiter)
	line, err = reader3.ReadBytes(delimiter)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Printf("Read contents(%d): %q\n", len(line), line)
	fmt.Printf("The number of unread bytes in the buffer: %d\n", reader3.Buffered())
	fmt.Println()

	// 示例6和示例7。
	fmt.Print("[ About contents leak ]\n\n")
	showContentsLeak(comment)
}

func showContentsLeak(comment string) {
	// 示例6。
	basicReader := strings.NewReader(comment)
	fmt.Printf("The size of basic reader: %d\n", basicReader.Size())

	size := len(comment)
	fmt.Printf("New a buffered reader with size %d ...\n", size)
	reader4 := bufio.NewReaderSize(basicReader, size)
	fmt.Println()

	peekNum := 7
	fmt.Printf("Peek %d bytes ...\n", peekNum)
	bytes, err := reader4.Peek(peekNum)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Printf("Peeked contents(%d): %q\n", len(bytes), bytes)
	fmt.Println()

	// 只要扩充一下之前拿到的字节切片bytes，
	// 就可以用它来读取甚至修改缓冲区中的后续内容。
	bytes = bytes[:cap(bytes)]
	fmt.Printf("The all of the contents in the buffer:\n%q\n", bytes)
	fmt.Println()

	blank := byte(' ')
	fmt.Println("Set blanks into the contents in the buffer ...")
	for _, i := range []int{55, 56, 57, 58, 66, 67, 68} {
		bytes[i] = blank
	}
	fmt.Println()

	peekNum = size
	fmt.Printf("Peek %d bytes ...\n", peekNum)
	bytes, err = reader4.Peek(peekNum)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Printf("Peeked contents(%d):\n%q\n", len(bytes), bytes)
	fmt.Println()

	// 示例7。
	// ReadSlice方法也存在相同的问题。
	delimiter := byte(',')
	fmt.Printf("Read slice with delimiter %q...\n", delimiter)
	line, err := reader4.ReadSlice(delimiter)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Printf("Read contents(%d): %q\n", len(line), line)
	fmt.Println()

	line = line[:cap(line)]
	fmt.Printf("The all of the contents in the buffer:\n%q\n", line)
	fmt.Println()

	underline := byte('_')
	fmt.Println("Set underlines into the contents in the buffer ...")
	for _, i := range []int{89, 92, 103} {
		line[i] = underline
	}
	fmt.Println()

	peekNum = size
	fmt.Printf("Peek %d bytes ...\n", peekNum)
	bytes, err = reader4.Peek(peekNum)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Printf("Peeked contents(%d): %q\n", len(bytes), bytes)
}

// TestBufIo_3
func TestBufIo_3(t *testing.T) {
	comment := "Writer implements buffering for an io.Writer object. " +
		"If an error occurs writing to a Writer, " +
		"no more data will be accepted and all subsequent writes, " +
		"and Flush, will return the error. After all data has been written, " +
		"the client should call the Flush method to guarantee all data " +
		"has been forwarded to the underlying io.Writer."
	basicWriter1 := &strings.Builder{}

	size := 300
	fmt.Printf("New a buffered writer with size %d ...\n", size)
	writer1 := bufio.NewWriterSize(basicWriter1, size)
	fmt.Println()

	// 示例1。
	begin, end := 0, 53
	fmt.Printf("Write %d bytes into the writer ...\n", end-begin)
	writer1.WriteString(comment[begin:end])
	fmt.Printf("The number of buffered bytes: %d\n", writer1.Buffered())
	fmt.Printf("The number of unused bytes in the buffer: %d\n",
		writer1.Available())
	fmt.Println("Flush the buffer in the writer ...")
	writer1.Flush()
	fmt.Printf("The number of buffered bytes: %d\n", writer1.Buffered())
	fmt.Printf("The number of unused bytes in the buffer: %d\n",
		writer1.Available())
	fmt.Println()

	// 示例2。
	begin, end = 0, 326
	fmt.Printf("Write %d bytes into the writer ...\n", end-begin)
	writer1.WriteString(comment[begin:end])
	fmt.Printf("The number of buffered bytes: %d\n", writer1.Buffered())
	fmt.Printf("The number of unused bytes in the buffer: %d\n",
		writer1.Available())
	fmt.Println("Flush the buffer in the writer ...")
	writer1.Flush()
	fmt.Println()

	// 示例3。
	basicWriter2 := &bytes.Buffer{}
	fmt.Printf("Reset the writer with a bytes buffer(an implementation of io.ReaderFrom) ...\n")
	writer1.Reset(basicWriter2)
	reader := strings.NewReader(comment)
	fmt.Println("Read data from the reader ...")
	writer1.ReadFrom(reader)
	fmt.Printf("The number of buffered bytes: %d\n", writer1.Buffered())
	fmt.Printf("The number of unused bytes in the buffer: %d\n",
		writer1.Available())
}
