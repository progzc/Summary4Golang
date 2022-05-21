package leetcode_0008_string_to_integer_atoi

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
)

// 0008.字符串转换整数 (atoi)
// https://leetcode-cn.com/problems/string-to-integer-atoi/

func Test_slice_1(t *testing.T) {
	s := []int{1, 2}
	s = append(s, 4, 5, 6)
	fmt.Printf("len=%d,cap=%d", len(s), cap(s))
}

func myAtoi(s string) int {
	// TODO
	return 0
	sync.Map{}
}
