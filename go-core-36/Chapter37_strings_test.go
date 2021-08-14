package go_core_36

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestStrings(t *testing.T) {
	// 示例1。
	var builder1 strings.Builder
	builder1.WriteString("A Builder is used to efficiently build a string using Write methods.")
	fmt.Printf("The first output(%d):\n%q\n", builder1.Len(), builder1.String())
	fmt.Println()
	builder1.WriteByte(' ')
	builder1.WriteString("It minimizes memory copying. The zero value is ready to use.")
	builder1.Write([]byte{'\n', '\n'})
	builder1.WriteString("Do not copy a non-zero Builder.")
	fmt.Printf("The second output(%d):\n\"%s\"\n", builder1.Len(), builder1.String())
	fmt.Println(builder1.Cap())
	fmt.Println("---------------------------------------------------------------")
	// 示例2。
	fmt.Println("Grow the builder ...")
	builder1.Grow(10)
	fmt.Printf("The length of contents in the builder is %d.\n", builder1.Len())
	fmt.Println(builder1.Cap())
	fmt.Println("---------------------------------------------------------------")

	// 示例3。
	fmt.Println("Reset the builder ...")
	builder1.Reset()
	fmt.Printf("The third output(%d):\n%q\n", builder1.Len(), builder1.String())
	fmt.Println(builder1.Cap())
}

func TestStrings2(t *testing.T) {
	// 示例1: 使用后就不要复制再使用了
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

	// 示例2：地址使用后可以复制再使用
	f2 := func(bp *strings.Builder) {
		(*bp).Grow(1) // 这里虽然不会引发panic，但不是并发安全的。
		builder4 := *bp
		//builder4.Grow(1) // 这里会引发panic。
		_ = builder4
	}
	f2(&builder1)

	// 零值可以被复制后再使用
	builder1.Reset()
	builder5 := builder1
	builder5.Grow(1) // 这里不会引发panic。
}

func TestStrings3(t *testing.T) {
	// 示例1。
	reader1 := strings.NewReader(
		"NewReader returns a new Reader reading from s. " +
			"It is similar to bytes.NewBufferString but more efficient and read-only.")
	fmt.Printf("The size of reader: %d\n", reader1.Size())
	fmt.Printf("The reading index in reader: %d\n",
		reader1.Size()-int64(reader1.Len()))

	buf1 := make([]byte, 47)
	n, _ := reader1.Read(buf1)
	fmt.Printf("%d bytes were read. (call Read)\n", n)
	fmt.Printf("The reading index in reader: %d\n",
		reader1.Size()-int64(reader1.Len()))
	fmt.Println()

	// 示例2。
	buf2 := make([]byte, 21)
	offset1 := int64(64)
	n, _ = reader1.ReadAt(buf2, offset1)
	fmt.Printf("%d bytes were read. (call ReadAt, offset: %d)\n", n, offset1)
	fmt.Printf("The reading index in reader: %d\n",
		reader1.Size()-int64(reader1.Len()))
	fmt.Println()

	// 示例3。
	offset2 := int64(17)
	expectedIndex := reader1.Size() - int64(reader1.Len()) + offset2
	fmt.Printf("Seek with offset %d and whence %d ...\n", offset2, io.SeekCurrent)
	readingIndex, _ := reader1.Seek(offset2, io.SeekCurrent)
	fmt.Printf("The reading index in reader: %d (returned by Seek)\n", readingIndex)
	fmt.Printf("The reading index in reader: %d (computed by me)\n", expectedIndex)

	n, _ = reader1.Read(buf2)
	fmt.Printf("%d bytes were read. (call Read)\n", n)
	fmt.Printf("The reading index in reader: %d\n",
		reader1.Size()-int64(reader1.Len()))
}
