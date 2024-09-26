package leetcode_0034_find_first_and_last_position_of_element_in_sorted_array

import "sort"

// 0034.在排序数组中查找元素的第一个和最后一个位置
// https://leetcode-cn.com/problems/find-first-and-last-position-of-element-in-sorted-array/

// searchRange_1 使用标准库
// 时间复杂度: O(log(n))
// 空间复杂度: O(1)
// 思路:
//
//	查找第一个元素	等价于	在数组中寻找第一个等于target的下标
//	查找最后一个元素	等于于	在数组中寻找最后一个等于target(即等于target+1)的下标,然后讲下标减一,并判断下标处的值是否等于target
func searchRange_1(nums []int, target int) []int {
	// sort.SearchInts采用快排,返回的是第一个等于target的索引
	li := sort.SearchInts(nums, target)
	if li == len(nums) || nums[li] != target {
		return []int{-1, -1}
	}
	ri := sort.SearchInts(nums, target+1) - 1
	return []int{li, ri}
}

// searchRange_2 使用4个变种
// 时间复杂度: O(log(n))
// 空间复杂度: O(1)
// 思路:
//
//	查找第一个元素	等价于	在数组中寻找第一个等于target的下标
//	查找最后一个元素	等于于	在数组中寻找最后一个等于target
func searchRange_2(nums []int, target int) []int {
	li, ri := -1, -1

	// 查找第一个等于target的下标
	low, high := 0, len(nums)-1
	for low <= high {
		mid := low + ((high - low) >> 1)
		if nums[mid] > target {
			high = mid - 1
		} else if nums[mid] < target {
			low = mid + 1
		} else {
			if (mid == 0) || (nums[mid-1] != target) { // 找到第一个与 target 相等的元素
				li = mid
				break
			}
			high = mid - 1
		}
	}

	if li == -1 {
		return []int{-1, -1}
	}

	// 查找最后一个等于target的下标
	low, high = 0, len(nums)-1
	for low <= high {
		mid := low + ((high - low) >> 1)
		if nums[mid] > target {
			high = mid - 1
		} else if nums[mid] < target {
			low = mid + 1
		} else {
			if (mid == len(nums)-1) || (nums[mid+1] != target) { // 找到最后一个与 target 相等的元素
				ri = mid
				break
			}
			low = mid + 1
		}
	}
	return []int{li, ri}
}
