package baidu_20241030_search_in_sorted_array

import (
	"fmt"
	"sort"
	"testing"
)

// 百度一面
// 在一个升序数组中找到第一次出现的目标数的位置，若找不到返回-1

func TestSearch(t *testing.T) {
	nums, target := []int{1, 2, 2, 3, 4, 5, 5}, 2
	fmt.Println(Search(nums, target)) // 1

	nums, target = []int{1, 2, 2, 3, 4, 5, 5}, 0
	fmt.Println(Search(nums, target)) // -1

	nums, target = []int{1}, 1
	fmt.Println(Search(nums, target)) // 0

	nums, target = []int{1}, 0
	fmt.Println(Search(nums, target)) // -1
}

func TestSearch_2(t *testing.T) {
	nums, target := []int{1, 2, 2, 3, 4, 5, 5}, 2
	fmt.Println(Search_2(nums, target)) // 1

	nums, target = []int{1, 2, 2, 3, 4, 5, 5}, 0
	fmt.Println(Search_2(nums, target)) // -1

	nums, target = []int{1}, 1
	fmt.Println(Search_2(nums, target)) // 0

	nums, target = []int{1}, 0
	fmt.Println(Search_2(nums, target)) // -1
}

// Search 变种二分
// 时间复杂度: O(log(n))
// 空间复杂度: O(1)
func Search(nums []int, target int) int {
	l, r := 0, len(nums)-1
	for l <= r {
		mid := l + (r-l)>>2
		if nums[mid] == target {
			if mid == 0 || nums[mid-1] != target {
				return mid
			} else {
				r = mid - 1
			}
		} else if nums[mid] > target {
			r = mid - 1
		} else {
			l = mid + 1
		}
	}
	return -1
}

// Search_2 直接调用 api
func Search_2(nums []int, target int) int {
	x := sort.SearchInts(nums, target)
	if x == len(nums) || nums[x] != target {
		return -1
	}
	return x
}
