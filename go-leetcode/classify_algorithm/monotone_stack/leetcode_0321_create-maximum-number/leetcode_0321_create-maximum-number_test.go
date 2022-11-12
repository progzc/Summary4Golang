package leetcode_0321_create_maximum_number

import (
	"fmt"
	"testing"
)

func Test_maxSubsequence(t *testing.T) {
	nums, k := []int{9, 1, 2, 5, 8, 3}, 3
	fmt.Println(maxSubsequence(nums, k))
}
