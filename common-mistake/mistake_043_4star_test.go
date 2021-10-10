package common_mistake

import (
	"bytes"
	"fmt"
	"testing"
)

// Slice 中数据的误用
func TestMistake_043(t *testing.T) {
	wrong043()
	right043()
	right043_1()
	right043_2()
	right043_3()
}

// 拼接的结果不是正确的 AAAAsuffix/BBBBBBBBB，因为 dir1、 dir2 两个 slice 引用的数据都是 path 的底层数组，
// 解决办法：
//    a. 重新分配新的 slice 并拷贝你需要的数据
//    b. 使用完整的 slice 表达式：input[low:high:max]，容量便调整为 max - low
//       关于go中的slice表达式：https://zhuanlan.zhihu.com/p/174742472
func wrong043() {
	path := []byte("AAAA/BBBBBBBBB")
	sepIndex := bytes.IndexByte(path, '/') // 4
	println(sepIndex)

	dir1 := path[:sepIndex]
	dir2 := path[sepIndex+1:]
	println("dir1: ", string(dir1)) // AAAA
	println("dir2: ", string(dir2)) // BBBBBBBBB

	dir1 = append(dir1, "suffix"...)
	println("current path: ", string(path)) // AAAAsuffixBBBB

	path = bytes.Join([][]byte{dir1, dir2}, []byte{'/'})
	println("dir1: ", string(dir1))     // AAAAsuffix
	println("dir2: ", string(dir2))     // uffixBBBB
	println("new path: ", string(path)) // AAAAsuffix/uffixBBBB    // 错误结果
}

func right043() {
	path := []byte("AAAA/BBBBBBBBB")
	sepIndex := bytes.IndexByte(path, '/') // 4
	// b. 使用完整的 slice 表达式：input[low:high:max]，容量便调整为 max - low
	dir1 := path[:sepIndex:sepIndex] // 此时 cap(dir1) 指定为4， 而不是先前的 16
	dir2 := path[sepIndex+1:]
	dir1 = append(dir1, "suffix"...)

	path = bytes.Join([][]byte{dir1, dir2}, []byte{'/'})
	println("dir1: ", string(dir1))     // AAAAsuffix
	println("dir2: ", string(dir2))     // BBBBBBBBB
	println("new path: ", string(path)) // AAAAsuffix/BBBBBBBBB
}

func right043_1() {
	numbers := [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	s := numbers[2:4:6]
	fmt.Println(s)              // [2, 3]
	fmt.Println(len(s), cap(s)) // 2 4

	s = append(s, []int{11, 12, 13}...)
	fmt.Println(s)              // [2 3 11 12 13]
	fmt.Println(len(s), cap(s)) // 5 8
	fmt.Println(numbers)        // [0 1 2 3 4 5 6 7 8 9]

	numbers[3] = 10
	fmt.Println(numbers) //[0 1 2 10 4 5 6 7 8 9]
	fmt.Println(s)       //[2 3 11 12 13]
}

func right043_2() {
	numbers := [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	s := numbers[1:4]
	fmt.Println(s)              // [1, 2, 3]
	fmt.Println(len(s), cap(s)) // 3 9

	s = numbers[1:4:5]
	fmt.Println(s)              // [1, 2, 3]
	fmt.Println(len(s), cap(s)) // 3 4
}

func right043_3() {
	numbers := [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Println(cap(numbers)) // 10
	s1 := numbers[1:4]
	fmt.Println(s1)      // [1, 2, 3]
	fmt.Println(cap(s1)) // 9
	s2 := numbers[1:4:5]
	fmt.Println(s2)      // [1, 2, 3]
	fmt.Println(cap(s2)) // 4
	// 当 slice 的输入操作数是一个 slice 时，结果 slice 的容量取决于输入操作数，而不是它的指向的底层 array
	s3 := s2[:]
	fmt.Println(s3)      // [1, 2, 3]
	fmt.Println(cap(s3)) // 4
}
