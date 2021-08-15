package go_core_36

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestBufio(t *testing.T) {
	comment := "Package bufio implements buffered I/O. " +
		"It wraps an io.Reader or io.Writer object, " +
		"creating another object (Reader or Writer) that " +
		"also implements the interface but provides buffering and " +
		"some help for textual I/O."
	basicReader := strings.NewReader(comment)
	fmt.Printf("The size of basic reader: %d\n", basicReader.Size())
	fmt.Println()

	// 示例1。
	fmt.Println("New a buffered reader ...")
	reader1 := bufio.NewReader(basicReader)
	fmt.Printf("The default size of buffered reader: %d\n", reader1.Size())
	// 此时reader1的缓冲区还没有被填充。
	fmt.Printf("The number of unread bytes in the buffer: %d\n", reader1.Buffered())
	fmt.Println()

	// 示例2。
	bytes, err := reader1.Peek(7)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Printf("Peeked contents(%d): %q\n", len(bytes), bytes)
	fmt.Printf("The number of unread bytes in the buffer: %d\n", reader1.Buffered())
	fmt.Println()

	// 示例3。
	buf1 := make([]byte, 7)
	n, err := reader1.Read(buf1)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Printf("Read contents(%d): %q\n", n, buf1)
	fmt.Printf("The number of unread bytes in the buffer: %d\n", reader1.Buffered())
	fmt.Println()

	// 示例4。
	fmt.Printf("Reset the basic reader (size: %d) ...\n", len(comment))
	basicReader.Reset(comment)
	fmt.Printf("Reset the buffered reader (size: %d) ...\n", reader1.Size())
	reader1.Reset(basicReader)
	peekNum := len(comment) + 1
	fmt.Printf("Peek %d bytes ...\n", peekNum)
	bytes, err = reader1.Peek(peekNum)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Printf("The number of peeked bytes: %d\n", len(bytes))
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
		fmt.Printf("error: %v\n", err)
	}
	fmt.Printf("The number of peeked bytes: %d\n", len(bytes))
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

func TestBufio2(t *testing.T) {
	comment := "Package bufio implements buffered I/O. " +
		"It wraps an io.Reader or io.Writer object, " +
		"creating another object (Reader or Writer) that " +
		"also implements the interface but provides buffering and " +
		"some help for textual I/O."
	basicReader := strings.NewReader(comment)
	fmt.Printf("The size of basic reader: %d\n", basicReader.Size())

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
	fmt.Printf("The number of unread bytes in the buffer: %d\n", reader1.Buffered())
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
	fmt.Printf("The number of unread bytes in the buffer: %d\n", reader1.Buffered())
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

func TestBufio3(t *testing.T) {
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
