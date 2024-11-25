package xiaomi_20241122_max_submatrix_lcci

import (
	"fmt"
	"math"
	"testing"
)

// 小米二面（端到端数据闭环工程师）

// 算法题1: 最大连续子数组的和（mid）
// 0053.最大子数组和
// https://leetcode.cn/problems/maximum-subarray

// 算法题2: 最大连续子矩阵的和（hard）🌟
// 面试题 17.24. 最大子矩阵
// https://leetcode.cn/problems/max-submatrix-lcci

func TestMaxSubArraySum(t *testing.T) {
	nums := []int{-2, 1, -3, 4, -1, 2, 1, -5, 4}
	fmt.Println(maxSubArraySum(nums)) // 6
}

func TestMaxSubMatrixSum(t *testing.T) {
	matrix := [][]int{
		{9, -8, 1, 3, -2},
		{-3, 7, 6, -2, 4},
		{6, -4, -4, 8, -7},
	}
	fmt.Println(maxSubMatrixSum(matrix)) // 19
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

// maxSubMatrixSum 最大连续子矩阵的和
// 思路：动态规划+前缀和
func maxSubMatrixSum(matrix [][]int) int {
	m := len(matrix)
	if m == 0 {
		return 0
	}
	n := len(matrix[0])
	var ans int
	for beginLine := 0; beginLine < m; beginLine++ {
		sum := make([]int, n)
		for i := beginLine; i < m; i++ {
			dp := math.MinInt32
			for j := 0; j < n; j++ {
				sum[j] += matrix[i][j]
				if dp > 0 {
					dp += sum[j]
				} else {
					dp = sum[j]
				}
				if dp > ans {
					ans = dp
				}
			}
		}
	}
	return ans
}
