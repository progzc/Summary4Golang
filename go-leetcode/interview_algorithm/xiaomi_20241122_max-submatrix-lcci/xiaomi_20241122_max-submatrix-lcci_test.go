package xiaomi_20241122_max_submatrix_lcci

import (
	"fmt"
	"testing"
)

// 小米二面
// 算法题1: 最大连续子数组的和（mid）
// 0053.最大子数组和
// https://leetcode.cn/problems/maximum-subarray

// 算法题2: 最大连续子矩阵的和（hard）

func TestMaxSubArraySum(t *testing.T) {
	nums := []int{-2, 1, -3, 4, -1, 2, 1, -5, 4}
	fmt.Println(maxSubArraySum(nums)) // 6
}

// maxSubArraySum 最大连续子数组的和
func maxSubArraySum(nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}
	if n == 1 {
		return nums[0]
	}
	dp := make([]int, n)
	dp[0] = nums[0]
	ans := dp[0]
	for i := 1; i < n; i++ {
		if dp[i-1] > 0 {
			dp[i] = dp[i-1] + nums[i]
		} else {
			dp[i] = nums[i]
		}
		if dp[i] > ans {
			ans = dp[i]
		}
	}
	return ans
}
