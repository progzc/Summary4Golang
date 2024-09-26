package leetcode_0033_search_in_rotated_sorted_array

// 0033.搜索旋转排序数组
// https://leetcode-cn.com/problems/search-in-rotated-sorted-array/

// search 二分搜索法
// 基本情况: nums升序排列,且互不相同
// 时间复杂度: O(log(n))
// 空间复杂度: O(1)
// 思路：每次分两半,必然有一半是有序的
func search(nums []int, target int) int {
	n := len(nums)
	left, right := 0, n-1
	for left <= right {
		mid := left + (right-left)/2
		if nums[mid] == target {
			return mid
		}
		if nums[mid] <= nums[right] { // 注意与右端点进行比较
			if nums[mid] < target && target <= nums[right] {
				left = mid + 1
			} else {
				right = mid - 1
			}
		} else {
			if nums[left] <= target && target < nums[mid] {
				right = mid - 1
			} else {
				left = mid + 1
			}
		}
	}
	return -1
}
