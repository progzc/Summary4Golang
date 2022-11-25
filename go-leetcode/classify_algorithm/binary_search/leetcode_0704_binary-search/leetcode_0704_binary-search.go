package leetcode_0704_binary_search

// 704. 二分查找
// https://leetcode.cn/problems/binary-search/

// search 二分查找
// 时间复杂度: O(log(n))
// 空间复杂度: O(1)
func search(nums []int, target int) int {
	low, high := 0, len(nums)-1
	for low <= high {
		mid := low + (high-low)>>2
		if nums[mid] == target {
			return mid
		} else if nums[mid] > target {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return -1
}
