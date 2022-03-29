package leetcode_0053_maximum_subarray

import "math"

// 0053.最大子数组和
// https://leetcode-cn.com/problems/maximum-subarray/

// maxSubArray 动态规划
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 思路：用f(i)代表以第i个数结尾的「连续子数组的最大和」，则f(i)=max{f(i−1)+nums[i],nums[i]}
func maxSubArray(nums []int) int {
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

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
