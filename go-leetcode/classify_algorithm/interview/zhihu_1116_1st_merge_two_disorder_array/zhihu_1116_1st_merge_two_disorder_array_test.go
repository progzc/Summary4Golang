package zhihu_1116_1st_merge_two_disorder_array

import (
	"fmt"
	"testing"
)

func Test_combine(t *testing.T) {
	nums1, nums2 := []int{4, 2, 3, 5, 0}, []int{9, 6, 4, 5, 2, 5}
	fmt.Println(combine(nums1, nums2))
}
