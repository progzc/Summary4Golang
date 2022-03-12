package leetcode_0035_search_insert_position

// 0035.搜索插入位置
// https://leetcode-cn.com/problems/search-insert-position/

// searchInsert 二分搜索
// 前提：升序+无重复
func searchInsert(nums []int, target int) int {
	low, high := 0, len(nums)-1
	for low <= high {
		mid := low + (high-low)/2
		if nums[mid] >= target {
			if (mid == 0) || (nums[mid-1] < target) {
				return mid
			}
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return len(nums)
}
