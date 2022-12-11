package leetcode_0265_paint_house_ii

import "math"

// 265. 粉刷房子 II
// https://leetcode.cn/problems/paint-house-ii/

// minCostII 动态规划
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路: 类似于[198. 打家劫舍]
//	状态: dp[i][j]表示粉刷第 0 号房子到第 i 号房子且第 i 号房子被粉刷成第 j 种颜色时的最小花费成本。
//	转移方程: dp[i][j] = min(dp[i-1][0],...,dp[i-1][j-1],dp[i-1][j+1]...,dp[i-1][n-1])+cost[i][j]
// 注意:
//	 此题代码可同样用于[256.粉刷房子]
func minCostII(costs [][]int) int {
	if len(costs) == 0 {
		return 0
	}

	n, m := len(costs), len(costs[0])
	dp := make([][]int, n)
	for i := 0; i < n; i++ {
		dp[i] = make([]int, m)
	}
	for j := 0; j < m; j++ {
		dp[0][j] = costs[0][j]
	}

	for i := 1; i < n; i++ {
		for j := 0; j < m; j++ {
			dp[i][j] = math.MaxInt32
			for k := 0; k < m; k++ {
				if k == j {
					continue
				}
				dp[i][j] = min(dp[i][j], costs[i][j]+dp[i-1][k])
			}
		}
	}

	ans := math.MaxInt32
	for j := 0; j < m; j++ {
		ans = min(ans, dp[n-1][j])
	}

	return ans
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
