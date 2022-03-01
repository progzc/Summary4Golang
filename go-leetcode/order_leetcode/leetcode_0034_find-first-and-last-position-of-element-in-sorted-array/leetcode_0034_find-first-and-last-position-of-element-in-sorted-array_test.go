package leetcode_0034_find_first_and_last_position_of_element_in_sorted_array

import "sort"

// searchRange1 二分法
// 时间复杂度 O(logn)
// 空间复杂度 O(1)
func searchRange1(nums []int, target int) []int {
	leftIdx := binarySearch(nums, target)
	rightIdx := binarySearch(nums, target+1) - 1
	if leftIdx <= rightIdx && rightIdx < len(nums) && nums[leftIdx] == target && nums[rightIdx] == target {
		return []int{leftIdx, rightIdx}
	}
	return []int{-1, -1}
}

// searchRange2 api调用
// 时间复杂度 O(logn)
// 空间复杂度 O(1)
func searchRange2(nums []int, target int) []int {
	leftIdx := sort.SearchInts(nums, target)
	if leftIdx == len(nums) || nums[leftIdx] != target {
		return []int{-1, -1}
	}
	rightIdx := sort.SearchInts(nums, target+1) - 1
	return []int{leftIdx, rightIdx}
}

func binarySearch(nums []int, target int) int {
	left, right, ans := 0, len(nums)-1, len(nums)
	for left <= right {
		mid := (left + right) / 2
		if nums[mid] >= target {
			right = mid - 1
			ans = mid
		} else {
			left = mid + 1
		}
	}
	return ans
}
