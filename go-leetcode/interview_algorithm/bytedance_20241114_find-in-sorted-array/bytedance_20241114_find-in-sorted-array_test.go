package bytedance_20241114_find_in_sorted_array

import (
	"fmt"
	"testing"
)

// 字节Es平台工程师（一面）
// 在重复的升序数组中找到第一个不小于目标数的序号

func TestFind(t *testing.T) {
	nums := []int{2, 2, 2, 2, 2, 2}
	target := 3
	fmt.Println(find(nums, target))
}

// find
func find(nums []int, target int) int {
	left, right := 0, len(nums)-1
	for left <= right {
		mid := left + (right-left)>>2
		if nums[mid] < target {
			left = mid + 1
		} else {
			if mid == 0 || nums[mid-1] < target {
				return mid
			} else {
				right = mid - 1
			}
		}
	}
	return -1
}
