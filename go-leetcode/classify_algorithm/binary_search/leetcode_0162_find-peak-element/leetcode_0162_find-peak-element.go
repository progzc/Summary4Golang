package leetcode_0162_find_peak_element

// 0162. 寻找峰值
// https://leetcode.cn/problems/find-peak-element/description/

// findPeakElement 二分法
// 时间复杂度: O(log(n))
// 空间复杂度: O(1)
func findPeakElement(nums []int) int {
	// todo
	return 0
}

// findPeakElement_2 寻找最大值
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 思路: 由于题目保证了 nums[i]≠nums[i+1]，那么数组 nums 中最大值两侧的元素一定严格小于最大值本身。
// 因此，最大值所在的位置就是一个可行的峰值位置。
func findPeakElement_2(nums []int) int {
	max, pos := nums[0], 0
	for i, v := range nums {
		if v > max {
			max = v
			pos = i
		}
	}
	return pos
}
