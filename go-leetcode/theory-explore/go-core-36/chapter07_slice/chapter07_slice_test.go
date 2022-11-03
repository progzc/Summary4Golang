package chapter07_slice

import (
	"fmt"
	"testing"
	"unsafe"
)

// TestSlice_1
// 关于切片的常见知识:
// (1) 切片不可比较,会直接编译报错
// (2) 切片扩容时,地址会不断变化,若没扩容,则地址不会变化
// (3) 切片的地址实际上是底层可变数组的首个元素的地址
func TestSlice_1(t *testing.T) {
	// 示例1。
	s1 := make([]int, 5)
	// s1; len: 5; cap: 5; value: [0 0 0 0 0]; address: 0xc0000d6030
	fmt.Printf("s1; len: %d; cap: %d; value: %d; address: %p\n", len(s1), cap(s1), s1, s1)
	s2 := make([]int, 5, 8)
	// s2; len: 5; cap: 8; value: [0 0 0 0 0]; address: 0xc0000c2100
	fmt.Printf("s2; len: %d; cap: %d; value: %d; address: %p\n", len(s2), cap(s2), s2, s2)
	fmt.Println()

	// 示例2。
	s3 := []int{1, 2, 3, 4, 5, 6, 7, 8}
	s4 := s3[3:6]
	// s3; len: 8; cap: 8; value: [1 2 3 4 5 6 7 8]; address: 0xc0000c2140
	fmt.Printf("s3; len: %d; cap: %d; value: %d; address: %p\n", len(s3), cap(s3), s3, s3)
	// s4; len: 3; cap: 5; value: [4 5 6]; address: 0xc0000c2158
	fmt.Printf("s4; len: %d; cap: %d; value: %d; address: %p\n", len(s4), cap(s4), s4, s4)
	fmt.Println()

	// 示例3。
	s5 := s4[:cap(s4)]
	// s5; len: 5; cap: 5; value: [4 5 6 7 8]; address: 0xc0000c2158
	fmt.Printf("s5; len: %d; cap: %d; value: %d; address: %p\n", len(s5), cap(s5), s5, s5)
	//fmt.Println(s4 == s5) 编译错误
	fmt.Println()

	s6 := []int{1, 2}
	// s6; len: 2; cap: 2; value: [1 2]; address: 0xc0000ac160
	fmt.Printf("s6; len: %d; cap: %d; value: %d; address: %p\n", len(s6), cap(s6), s6, s6)
	s6 = append(s6, 3)
	// s6; len: 3; cap: 4; value: [1 2 3]; address: 0xc0000aa080
	fmt.Printf("s6; len: %d; cap: %d; value: %d; address: %p\n", len(s6), cap(s6), s6, s6)
	s6 = append(s6, []int{3, 4, 5, 6, 7, 8, 9, 10, 11, 12}...)
	// s6; len: 13; cap: 14; value: [1 2 3 3 4 5 6 7 8 9 10 11 12]; address: 0xc0000e2000
	fmt.Printf("s6; len: %d; cap: %d; value: %d; address: %p\n", len(s6), cap(s6), s6, s6)
}

// 关于切片的内存分配:
// 源码见 D:/Program Files/go/go1.16.8/src/runtime/slice.go:125中的growslice函数
// (1) 初始化切片时,若未指定容量,则容量等于长度
// (2) 容量小于1024,扩容时按照1.5倍(实际考虑到内存对齐,会有细微差异)
// (3) 容量大于1024,扩容按照1.25倍(实际考虑到内存对齐,会有细微差异)
// (4) 切片不会缩容
func TestSlice_2(t *testing.T) {
	// 示例1。
	s6 := make([]int, 0)
	fmt.Printf("The capacity of s6: %d\n", cap(s6))
	for i := 1; i <= 5; i++ {
		s6 = append(s6, i)
		fmt.Printf("s6(%d): len: %d, cap: %d\n", i, len(s6), cap(s6))
	}
	fmt.Println()

	// 示例2。
	s7 := make([]int, 1024)
	fmt.Printf("The capacity of s7: %d\n", cap(s7))
	s7e1 := append(s7, make([]int, 200)...)
	fmt.Printf("s7e1: len: %d, cap: %d\n", len(s7e1), cap(s7e1))
	s7e2 := append(s7, make([]int, 400)...)
	fmt.Printf("s7e2: len: %d, cap: %d\n", len(s7e2), cap(s7e2))
	s7e3 := append(s7, make([]int, 600)...)
	fmt.Printf("s7e3: len: %d, cap: %d\n", len(s7e3), cap(s7e3))
	fmt.Println()

	// 示例3。
	s8 := make([]int, 10)
	fmt.Printf("The capacity of s8: %d\n", cap(s8))
	s8a := append(s8, make([]int, 11)...)
	fmt.Printf("s8a: len: %d, cap: %d\n", len(s8a), cap(s8a))
	s8b := append(s8a, make([]int, 23)...)
	fmt.Printf("s8b: len: %d, cap: %d\n", len(s8b), cap(s8b))
	s8c := append(s8b, make([]int, 45)...)
	fmt.Printf("s8c: len: %d, cap: %d\n", len(s8c), cap(s8c))
}

// TestSlice_3 一道常见面试题
// 来源: https://eddycjy.com/posts/go/slice-discuss/
// 要点：
// (1) 理解Slice的底层结构(如下): slice真正存储数据的地方，是一个数组。slice 的结构中存储的是指向所引用的数组指针地址
// (2) 函数都是值传递,需要复制（切片的本质也是值,而不是指针）
type slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}

func TestSlice_3(t *testing.T) {
	s1 := make([]int, 0, 10)
	var appenFunc = func(s []int) {
		s = append(s, 10, 20, 30)
		// value: [], address: 0xc000016230, pointer address: 0xc000004090
		fmt.Printf("value: %d, address: %p, pointer address: %p\n", s1, s1, &s1)
	}
	// value: [], address: 0xc000016230, pointer address: 0xc000004090
	fmt.Printf("value: %d, address: %p, pointer address: %p\n", s1, s1, &s1)
	appenFunc(s1)
	fmt.Println(s1)      // []
	fmt.Println(s1[:10]) // [10 20 30 0 0 0 0 0 0 0]
	fmt.Println()

	// -------------------------------------------------------------
	s2 := make([]int, 0, 10)
	var appenFunc2 = func(s *[]int) {
		*s = append(*s, 10, 20, 30)
		// value: &[10 20 30], address: 0xc000004138, pointer address: 0xc000006038
		fmt.Printf("value: %d, address: %p, pointer address: %p\n", s, s, &s)
	}
	// value: [], address: 0xc000016280, pointer address: 0xc000004138
	fmt.Printf("value: %d, address: %p, pointer address: %p\n", s2, s2, &s2)
	appenFunc2(&s2)
	fmt.Println(s2)      // [10 20 30]
	fmt.Println(s2[:10]) // [10 20 30 0 0 0 0 0 0 0]
}

func Test_slice(t *testing.T) {
	s := []int{5}
	s = append(s, 7)
	s = append(s, 9)
	x := append(s, 11)
	y := append(s, 12)
	fmt.Println(s, x, y) // [5 7 9] [5 7 9 12] [5 7 9 12]
}

func Test_slice_2(t *testing.T) {
	s := []int{5}
	s = append(s, 7)
	s = append(s, 9)
	x := append(s, 11)
	y := append(s, 12, 13)
	fmt.Println(s, x, y) // [5 7 9] [5 7 9 11] [5 7 9 12 13]
}
