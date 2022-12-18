package leetcode_0120_triangle

import "math"

// 120. 三角形最小路径和
// https://leetcode.cn/problems/triangle/

// minimumTotal dfs (超时)
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func minimumTotal(triangle [][]int) int {
	var (
		dfs func(layer, i, sum int)
		n   = len(triangle)
		ans = math.MaxInt32
	)
	// layer 表示层数, i 表示下标, sum 总和
	dfs = func(layer, i, sum int) {
		if layer == n {
			ans = min(ans, sum)
			return
		}
		if i >= len(triangle[layer]) {
			return
		}
		dfs(layer+1, i, sum+triangle[layer][i])
		dfs(layer+1, i+1, sum+triangle[layer][i])
	}
	dfs(0, 0, 0)
	return ans
}

// minimumTotal_2 动态规划（二维）
// 时间复杂度: O(n^2)
// 空间复杂度: O(n^2)
func minimumTotal_2(triangle [][]int) int {
	n := len(triangle)
	dp := make([][]int, n)
	for i := 0; i < n; i++ {
		dp[i] = make([]int, n)
	}
	dp[0][0] = triangle[0][0]
	for i := 1; i < n; i++ {
		dp[i][0] = dp[i-1][0] + triangle[i][0]
	}

	for i := 1; i < n; i++ {
		for j := 1; j <= i; j++ {
			if j < len(triangle[i-1]) {
				dp[i][j] = min(dp[i-1][j], dp[i-1][j-1]) + triangle[i][j]
			} else {
				dp[i][j] = dp[i-1][j-1] + triangle[i][j]
			}
		}
	}

	ans := math.MaxInt32
	for _, sum := range dp[n-1] {
		ans = min(ans, sum)
	}
	return ans
}

// minimumTotal_3 动态规划（一维）
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func minimumTotal_3(triangle [][]int) int {
	n := len(triangle)
	dp := make([]int, n)
	dp[0] = triangle[0][0]

	for i := 1; i < n; i++ {
		for j := i; j >= 1; j-- {
			if j < len(triangle[i-1]) {
				dp[j] = min(dp[j], dp[j-1]) + triangle[i][j]
			} else {
				dp[j] = dp[j-1] + triangle[i][j]
			}
		}
		dp[0] = dp[0] + triangle[i][0]
	}

	ans := math.MaxInt32
	for _, sum := range dp {
		ans = min(ans, sum)
	}
	return ans
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
