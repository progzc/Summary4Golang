package go_bible

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"testing"
	"time"
)

func square(n int) int     { return n * n }
func negative(n int) int   { return -n }
func product(m, n int) int { return m * n }
func TestFunc(t *testing.T) {
	f := square
	fmt.Println(f(3)) // "9"

	f = negative
	fmt.Println(f(3))     // "-3"
	fmt.Printf("%T\n", f) // "func(int) int"

	// compile error: can't assign func(int, int) int to func(int) int
	//f = product
}

func TestFunc2(t *testing.T) {
	var f func(int) int
	// 此处f的值为nil, 会引起panic错误
	f(3)
}

func TestFunc3(t *testing.T) {
	var f func(int) int
	if f != nil {
		f(3)
	}
}

func add1(r rune) rune {
	return r + 1
}
func TestFunc4(t *testing.T) {
	fmt.Println(strings.Map(add1, "HAL-9000")) // "IBM.:111"
	fmt.Println(strings.Map(add1, "VMS"))      // "WNT"
	fmt.Println(strings.Map(add1, "Admix"))    // "Benjy"
}

func TestAnonymityFunc(t *testing.T) {
	// 使用匿名函数
	fmt.Println(strings.Map(func(r rune) rune { return r + 1 }, "HAL-9000")) // "IBM.:111"
}

func squares() func() int {
	var x int
	return func() int {
		x++
		return x * x
	}
}
func TestAnonymityFunc2(t *testing.T) {
	f := squares()
	fmt.Println(f()) // "1"
	fmt.Println(f()) // "4"
	fmt.Println(f()) // "9"
	fmt.Println(f()) // "16"
}

func TestAnonymityFunc3(t *testing.T) {
	var prereqs = map[string][]string{
		"algorithms": {"data structures"},
		"calculus":   {"linear algebra"},
		"compilers": {
			"data structures",
			"formal languages",
			"computer organization",
		},
		"data structures":       {"discrete math"},
		"databases":             {"data structures"},
		"discrete math":         {"intro to programming"},
		"formal languages":      {"discrete math"},
		"networks":              {"operating systems"},
		"operating systems":     {"data structures", "computer organization"},
		"programming languages": {"data structures", "computer organization"},
	}
	for i, course := range topoSort(prereqs) {
		fmt.Printf("%d:\t%s\n", i+1, course)
	}
}
func topoSort(m map[string][]string) []string {
	var order []string
	seen := make(map[string]bool)
	// 当匿名函数需要被递归调用时,我们必须首先声明一个变量（在上面的例子中,我们首先声明了 visitAll）, 再将匿名函数赋值给这个变量。
	// 如果不分成两步，函数字面量无法与visitAll绑定，我们也无法递归调用该匿名函数。
	var visitAll func(items []string)
	visitAll = func(items []string) {
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				visitAll(m[item])
				order = append(order, item)
			}
		}
	}
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	visitAll(keys)
	return order
}

func trace(msg string) func() {
	start := time.Now()
	log.Printf("enter %s", msg)
	return func() {
		log.Printf("exit %s (%s)", msg, time.Since(start))
	}
}
func bigSlowOperation() {
	defer trace("bigSlowOperation")() // don't forget the extra parentheses
	// ...lots of work…
	time.Sleep(10 * time.Second) // simulate slow operation by sleeping
}
func TestDefer(t *testing.T) {
	bigSlowOperation()
}

func double(x int) (result int) {
	defer func() { fmt.Printf("double(%d) = %d\n", x, result) }()
	return x + x
}
func TestDefer2(t *testing.T) {
	_ = double(4)
}

func triple(x int) (result int) {
	defer func() { result += x }()
	return double(x)
}
func TestDefer3(t *testing.T) {
	fmt.Println(triple(4))
}
