package chapter07_slice

import (
	"fmt"
	"testing"
)

func Test_interview_slice_1(t *testing.T) {
	s := []int{5}
	s = append(s, 7)
	s = append(s, 9)
	x := append(s, 11)
	y := append(s, 12)
	fmt.Println(s, x, y) // [5 7 9] [5 7 9 12] [5 7 9 12]
}

func Test_interview_slice_2(t *testing.T) {
	s := []int{5}
	s = append(s, 7)
	s = append(s, 9)
	x := append(s, 11)
	y := append(s, 12, 13)
	fmt.Println(s, x, y) // [5 7 9] [5 7 9 11] [5 7 9 12 13]
}

func Test_interview_slice_3(t *testing.T) {
	s1 := make([]int, 0, 10)
	var appendFunc = func(s []int) {
		s = append(s, 10, 20, 30)
		fmt.Println(s) // [10 20 30]
	}
	fmt.Println(s1) // []
	appendFunc(s1)
	fmt.Println(s1)      // []
	fmt.Println(s1[:10]) // [10 20 30 0 0 0 0 0 0 0]
	fmt.Println(s1[:])   // []
}

func Test_interview_slice_4(t *testing.T) {
	s1 := make([]int, 0, 10)
	var appendFunc = func(s *[]int) {
		*s = append(*s, 10, 20, 30)
		fmt.Println(s) // &[10 20 30]
	}
	fmt.Println(s1) // []
	appendFunc(&s1)
	fmt.Println(s1)      // [10 20 30]
	fmt.Println(s1[:10]) // [10 20 30 0 0 0 0 0 0 0]
	fmt.Println(s1[:])   // [10 20 30]
}

func Test_interview_slice_5(t *testing.T) {
	s1 := []int{1, 2, 3, 4}
	s2 := make([]int, 2)
	l := copy(s2, s1)
	fmt.Println(s2) // [1 2]
	fmt.Println(l)  // 2
}

func Test_interview_slice_6(t *testing.T) {
	x := make([]int, 0, 10)
	x = append(x, 1, 2, 3)
	y := append(x, 4)
	z := append(x, 5)
	fmt.Println(x) // [1 2 3]
	fmt.Println(y) // [1 2 3 5]
	fmt.Println(z) // [1 2 3 5]
}

func Test_interview_slice_7(t *testing.T) {
	a := []int{1, 2, 3}     // 长度为3，容量为3
	b := make([]int, 1, 10) // 长度为1，容量为10
	modify(a, b)
	fmt.Println("out a =", a) // out a = [1 2 3]
	fmt.Println("out b =", b) // out b = [3]
}

func modify(a, b []int) {
	a = append(a, 4)            // 引发扩容，此时返回的a是一个新的切片
	b = append(b, 2)            // 没有引发扩容，仍然是原切片
	a[0] = 3                    //改变a切片元素
	b[0] = 3                    //改变b切片元素
	fmt.Println("inner a =", a) // inner a = [3 2 3 4]
	fmt.Println("inner b =", b) // inner b = [3 2]
}

func Test_interview_array_1(t *testing.T) {
	nums := [4]int{1, 2, 3, 4}
	modifyArray(nums)
	fmt.Println(nums) // [1 2 3 4]
}

func modifyArray(nums [4]int) {
	nums[0] = 10
}

func Test_interview_array_2(t *testing.T) {
	arr := [5]int{1, 2, 3, 4, 5}
	m := make(map[int]int, 5)
	for i, v := range arr {
		m[i] = v
	}
	fmt.Println(m)
}
