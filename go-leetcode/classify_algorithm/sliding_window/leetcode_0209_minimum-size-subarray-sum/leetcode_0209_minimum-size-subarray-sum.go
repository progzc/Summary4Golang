package leetcode_0209_minimum_size_subarray_sum

import "math"

// 0209. 长度最小的子数组
// https://leetcode.cn/problems/minimum-size-subarray-sum/

// minSubArrayLen 滑动窗口
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func minSubArrayLen(target int, nums []int) int {
	n := len(nums)
	ans, sum := math.MaxInt32, 0
	for l, r := 0, 0; r < n; r++ {
		sum += nums[r]
		for sum >= target {
			ans = min(ans, r-l+1)
			sum -= nums[l]
			l++
		}
	}
	if ans == math.MaxInt32 {
		return 0
	}
	return ans
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
