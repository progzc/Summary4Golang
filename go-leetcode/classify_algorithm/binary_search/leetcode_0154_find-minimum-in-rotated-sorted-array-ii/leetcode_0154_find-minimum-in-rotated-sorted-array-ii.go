package leetcode_0154_find_minimum_in_rotated_sorted_array_ii

// 0154.寻找旋转排序数组中的最小值 II
// https://leetcode-cn.com/problems/find-minimum-in-rotated-sorted-array-ii/

// findMin
// 基本情况: 升序+有相同+旋转
// 思路: 画曲线图 + 思考特殊情况
// 时间复杂度: O(log(n))
// 空间复杂度: O(1)
func findMin(nums []int) int {
	left, right := 0, len(nums)-1
	// 这里不取等号(因为下面赋值是right = mid),否则如果最小值恰好在中间时会出错
	// 例如: [2,4,0,2,2]
	for left < right {
		mid := left + (right-left)/2
		// 由于互不重复,不可能出现 nums[mid] == nums[right] 的情况
		if nums[mid] < nums[right] {
			right = mid
		} else if nums[mid] > nums[right] {
			left = mid + 1
		} else {
			right--
		}
	}
	return nums[left]
}
