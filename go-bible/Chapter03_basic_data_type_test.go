package go_bible

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
	"unicode/utf8"
)

func TestString(t *testing.T) {
	s := "中国"
	fmt.Println(len(s), s[0], s[1], s[2])
	fmt.Println(s[0:6])
}

func TestString2(t *testing.T) {
	s := "Hello, 世界"
	fmt.Println(len(s))                    // "13"
	fmt.Println(utf8.RuneCountInString(s)) // "9"

	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		fmt.Printf("%d\t%c\n", i, r)
		i += size
	}
}

func TestString3(t *testing.T) {
	s := "Hello, 世界"
	fmt.Println(len(s))                    // "13"
	fmt.Println(utf8.RuneCountInString(s)) // "9"

	for i, r := range "Hello, 世界" {
		fmt.Printf("%d\t%q\t%d\n", i, r, r)
	}
}

func TestString4(t *testing.T) {
	s := "Hello, 世界"
	n := 0
	for range s {
		n++
	}
	fmt.Println(n)
	fmt.Println(utf8.RuneCountInString(s))
}

func TestString5(t *testing.T) {
	fmt.Println(basename("a/b/c.go_knowledge")) // "c"
	fmt.Println(basename("c.d.go_knowledge"))   // "c.d"
	fmt.Println(basename("abc"))                // "abc"
}

func basename(s string) string {
	slash := strings.LastIndex(s, "/") // -1 if "/" not found
	s = s[slash+1:]
	if dot := strings.LastIndex(s, "."); dot >= 0 {
		s = s[:dot]
	}
	return s
}

func TestString6(t *testing.T) {
	fmt.Println(comma("12345")) // "12,345"
}

func comma(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	return comma(s[:n-3]) + "," + s[n-3:]
}

func TestBuffer(t *testing.T) {
	fmt.Println(intsToString([]int{1, 2, 3}))
}
func intsToString(values []int) string {
	var buf bytes.Buffer
	buf.WriteByte('[')
	for i, v := range values {
		if i > 0 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(&buf, "%d", v)
	}
	buf.WriteRune(rune(']'))
	return buf.String()
}

func TestStringConvert(t *testing.T) {
	x := 123
	y := fmt.Sprintf("%d", x)
	fmt.Println(y, strconv.Itoa(x)) // "123 123"
	fmt.Println(strconv.FormatInt(int64(x), 2))
}

func TestStringConvert2(t *testing.T) {
	x, _ := strconv.Atoi("123")
	y, _ := strconv.ParseInt("123", 10, 64)
	fmt.Println(x, y)
}

const noDelay time.Duration = 0
const timeout = 5 * time.Minute

func TestConst(t *testing.T) {
	fmt.Printf("%T %[1]v\n", noDelay)     // "time.Duration 0"
	fmt.Printf("%T %[1]v\n", timeout)     // "time.Duration 5m0s"
	fmt.Printf("%T %[1]v\n", time.Minute) // "time.Duration 1m0s"
}

const (
	a = 1
	b
	c = 2
	d
)

func TestConst2(t *testing.T) {
	fmt.Println(a, b, c, d) // "1 1 2 2"
}

type Weekday int

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

func TestConst3(t *testing.T) {
	fmt.Printf("%d\t%d\n", Sunday, Monday)
}

type Flags uint

const (
	FlagUp           Flags = 1 << iota // is up
	FlagBroadcast                      // supports broadcast access capability
	FlagLoopback                       // is a loopback interface
	FlagPointToPoint                   // belongs to a point-to-point link
	FlagMulticast                      // supports multicast access capability
)

func TestConst4(t *testing.T) {
	fmt.Printf("%d\t%d\t%d\n", FlagUp, FlagBroadcast, FlagLoopback)
}

const (
	_   = 1 << (10 * iota)
	KiB // 1024
	MiB // 1048576
	GiB // 1073741824
	TiB // 1099511627776             (exceeds 1 << 32)
	PiB // 1125899906842624
	EiB // 1152921504606846976
	ZiB // 1180591620717411303424    (exceeds 1 << 64)
	YiB // 1208925819614629174706176
)

func TestConst5(t *testing.T) {
	fmt.Printf("%d\t%d\t%d\n", KiB, MiB, GiB)
}
