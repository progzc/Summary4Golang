package leetcode_0035_search_insert_position

// searchInsert 二分法
// 时间复杂度 O(logn)
// 空间复杂度 O(1)
func searchInsert(nums []int, target int) int {
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
