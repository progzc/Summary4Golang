package bytedance_20241029_merge_two_sorted_array

import (
	"fmt"
	"testing"
)

func TestMerge(t *testing.T) {
	nums1 := []int{1, 2, 2, 2, 3, 5}
	nums2 := []int{2, 3, 3, 6}
	fmt.Println(merge(nums1, nums2)) // [1 2 2 2 3 5 2 3 3 6]

	nums1 = []int{1, 2, 2, 2, 3, 5}
	nums2 = []int{}
	fmt.Println(merge(nums1, nums2)) // [1 2 2 2 3 5]
}

// 字节二面-合并2个排序数组
func merge(nums1, nums2 []int) []int {
	m, n := len(nums1), len(nums2)
	i, j := 0, 0
	var ans []int
	if i < m && j < n {
		if nums1[i] <= nums2[j] {
			ans = append(ans, nums1[i])
			i++
		} else {
			ans = append(ans, nums2[j])
			j++
		}
	}

	if i < m {
		ans = append(ans, nums1[i:]...)
	}

	if j < n {
		ans = append(ans, nums2[j:]...)
	}
	return ans
}
