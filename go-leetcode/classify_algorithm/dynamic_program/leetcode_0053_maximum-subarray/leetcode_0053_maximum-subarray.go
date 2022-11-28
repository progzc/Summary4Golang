package leetcode_0053_maximum_subarray

import "math"

// 0053.最大子数组和
// https://leetcode-cn.com/problems/maximum-subarray/

// maxSubArray 动态规划
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路：
//	状态: dp[i] 代表以元素 nums[i] 为结尾的连续子数组最大和。
//	转移方程: 若 dp[i-1] ≤ 0 ，说明 dp[i-1] 对 dp[i] 产生负贡献，即 dp[i-1]+nums[i] 还不如 nums[i] 本身大.
//		当dp[i-1]>0时, dp[i] = dp[i-1] + nums[i];
//		当dp[i-1]≤0时, dp[i] = nums[i];
//	初始状态: dp[0] = nums[0], 即以 nums[0] 结尾的连续子数组最大和为 nums[0].
func maxSubArray(nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}
	if n == 1 {
		return nums[0]
	}

	dp := make([]int, n)
	dp[0] = nums[0]
	ans := math.MinInt32
	ans = max(ans, dp[0])
	for i := 1; i < n; i++ {
		if dp[i-1] > 0 {
			dp[i] = dp[i-1] + nums[i]
		} else {
			dp[i] = nums[i]
		}
		ans = max(ans, dp[i])
	}
	return ans
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// maxSubArray_3 动态规划
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func maxSubArray_3(nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}

	pre := nums[0]
	ans := nums[0]
	for i := 1; i < n; i++ {
		var cur int
		if pre < 0 {
			cur = nums[i]
		} else {
			cur = pre + nums[i]
		}
		ans = max(ans, cur)
		pre = cur
	}
	return ans
}

// maxSubArray_2 动态规划
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 思路：
//	状态: dp[i] 代表以元素 nums[i] 为结尾的连续子数组最大和。
//	转移方程: 若 dp[i-1] ≤ 0 ，说明 dp[i-1] 对 dp[i] 产生负贡献，即 dp[i-1]+nums[i] 还不如 nums[i] 本身大.
//		当dp[i-1]>0时, dp[i] = dp[i-1] + nums[i];
//		当dp[i-1]≤0时, dp[i] = nums[i];
//	初始状态: dp[0] = nums[0], 即以 nums[0] 结尾的连续子数组最大和为 nums[0].
func maxSubArray_2(nums []int) int {
	ans, pre := math.MinInt32, math.MinInt32
	for _, num := range nums {
		if pre+num < num {
			pre = num
		} else {
			pre = num + pre
		}
		ans = max(ans, pre)
	}
	return ans
}
