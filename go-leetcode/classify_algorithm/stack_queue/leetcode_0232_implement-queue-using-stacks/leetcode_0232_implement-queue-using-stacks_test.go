package leetcode_0232_implement_queue_using_stacks

import (
	"fmt"
	"testing"
)

func Test_slice(t *testing.T) {
	s1 := []int{1}
	s1 = s1[1:] // 不会出现越界错误
	fmt.Println(s1)

	s2 := []int{1, 2}
	s2 = s2[2:] // 不会出现越界错误
	fmt.Println(s2)

}
