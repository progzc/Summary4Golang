package leetcode_0279_perfect_squares

import "math"

// 0279.完全平方数🌟
// https://leetcode-cn.com/problems/perfect-squares/

// numSquares 动态规划(完全背包问题)
// 时间复杂度: O(n*sqrt(n))
// 空间复杂度: O(n)
// 思路：转化为完全背包
// 特点：完全背包的最值问题
func numSquares(n int) int {
	// dp[i]: 和为i的完全平方数的最小数量
	dp := make([]int, n+1)
	for i := 0; i < len(dp); i++ {
		dp[i] = math.MaxInt32
	}
	// 和为0的完全平方数的最小数量为0
	dp[0] = 0
	for num := 1; num*num <= n; num++ {
		for j := 1; j <= n; j++ {
			if j >= num*num {
				dp[j] = min(dp[j], dp[j-num*num]+1)
			}
		}
	}
	return dp[n]
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// TODO 优化解法
