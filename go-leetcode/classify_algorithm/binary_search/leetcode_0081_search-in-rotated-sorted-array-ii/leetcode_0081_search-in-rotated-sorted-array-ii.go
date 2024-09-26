package leetcode_0081_search_in_rotated_sorted_array_ii

// 0081.搜索旋转排序数组 II
// https://leetcode-cn.com/problems/search-in-rotated-sorted-array-ii/

// search 二分搜索法
// 基本情况: nums升序排列,且有相同, 然后旋转
// 时间复杂度: O(log(n))（最坏情况下是O(n)）
// 空间复杂度: O(1)
// 思路：每次分两半,必然有一半是有序的；但有一种特例会出现导致无法判断哪部分有序：nums[left] == nums[mid] == nums[right]
func search(nums []int, target int) bool {
	n := len(nums)
	left, right := 0, n-1
	for left <= right {
		mid := left + (right-left)/2
		if nums[mid] == target {
			return true
		}
		// 如果nums[left] <= nums[mid]，则必然前半部分是有序的;否则后半部分有序
		// 但有一种特例会出现导致无法判断哪部分有序：nums[left] == nums[mid] == nums[right]
		if nums[left] == nums[mid] && nums[mid] == nums[right] {
			left++
			right--
			continue
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
	return false
}
