package leetcode_0162_find_peak_element

// 0162. 寻找峰值
// https://leetcode.cn/problems/find-peak-element/description/

// findPeakElement 二分法
// 时间复杂度: O(log(n))
// 空间复杂度: O(1)
// 思路：在题目描述中出现了 nums[-1] = nums[n] = -∞，这就代表着 只要数组中存在一个元素比相邻元素大，那么沿着它一定可以找到一个峰值。
func findPeakElement(nums []int) int {
	left, right := 0, len(nums)-1
	for left < right { // 不能取等于，不然就会出现nums[mid+1]越界的情况
		mid := left + (right-left)/2
		if nums[mid] > nums[mid+1] { // 不会出现等于的情况，否则就无法进行搜索了。
			right = mid
		} else {
			left = mid + 1
		}
	}
	return left
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
