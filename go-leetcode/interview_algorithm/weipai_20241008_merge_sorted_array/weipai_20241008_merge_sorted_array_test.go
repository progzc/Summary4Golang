package weipai_20241008_merge_sorted_array

import (
	"fmt"
	"testing"
)

func TestMerges(t *testing.T) {
	nums := [][]int{{1, 2, 3}, {2, 3, 4}, {2, 4, 5}, {4, 5, 6}}
	//m := len(nums)
	//fmt.Println(nums[:m/2]) // [[1 2 3]]
	//fmt.Println(nums[m/2:]) // [[2 3 4] [2 4 5]]
	fmt.Println(Merges(nums)) // [1 2 2 2 3 3 4 4 4 5 5 6]
}

func TestMerge(t *testing.T) {
	nums1 := []int{1, 2, 3}
	nums2 := []int{2, 2, 3}
	fmt.Println(merge(nums1, nums2)) // [1 2 2 2 3 3]
}

func Merges(nums [][]int) []int {
	var ans []int
	m := len(nums)
	if m == 0 {
		return ans
	}
	if m == 1 {
		ans = nums[0]
		return ans
	}
	return merge(Merges(nums[:m/2]), Merges(nums[m/2:]))
}

func merge(nums1, nums2 []int) []int {
	var ans []int
	l1, l2 := len(nums1), len(nums2)
	i, j := 0, 0
	for i < l1 && j < l2 {
		if nums1[i] <= nums2[j] {
			ans = append(ans, nums1[i])
			i++
		} else {
			ans = append(ans, nums2[j])
			j++
		}
	}
	if i < l1 {
		ans = append(ans, nums1[i:]...)
	}
	if j < l2 {
		ans = append(ans, nums2[j:]...)
	}
	return ans
}
