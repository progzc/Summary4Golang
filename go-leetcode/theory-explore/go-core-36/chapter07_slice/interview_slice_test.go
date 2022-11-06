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
