package leetcode_0153_find_minimum_in_rotated_sorted_array

// 0153.寻找旋转排序数组中的最小值
// https://leetcode-cn.com/problems/find-minimum-in-rotated-sorted-array/

// findMin
// 基本情况: 升序+互不相同+旋转
// 思路: 画曲线图
// 时间复杂度: O(log(n))
// 空间复杂度: O(1)
func findMin(nums []int) int {
	left, right := 0, len(nums)-1
	// 这里不取等号(因为下面赋值是right = mid),否则如果最小值恰好在中间时会导致进入死循环
	// 例如: [3,4,0,1,2]
	for left < right {
		mid := left + (right-left)/2
		if nums[mid] <= nums[right] { // 必须和nums[right]进行比较
			right = mid
		} else {
			left = mid + 1
		}
	}
	return nums[left]
}
