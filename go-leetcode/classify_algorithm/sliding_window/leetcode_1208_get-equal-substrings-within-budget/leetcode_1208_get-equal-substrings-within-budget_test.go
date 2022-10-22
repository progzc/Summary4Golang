package leetcode_1208_get_equal_substrings_within_budget

import (
	"fmt"
	"testing"
)

func Test_byte(t *testing.T) {
	ss, tt := "abcd", "bcdf"
	fmt.Println(ss[0], tt[0])
	fmt.Println(ss[0]-'a', tt[0]-'a')
	fmt.Println(ss[0]-tt[0], int(ss[0])-int(tt[0]))
}

func Test_diff(t *testing.T) {
	ss, tt := "pxezla", "loewbi"
	diff := make([]int, len(ss))
	for i := 0; i < len(ss); i++ {
		diff[i] = abs(int(ss[i]) - int(tt[i]))
	}
	fmt.Println(diff)
}
