package leetcode_0256_paint_house

import "math"

// 0256. 粉刷房子
// https://leetcode.cn/problems/paint-house/

// minCost dfs（超时）
// 时间复杂度: O(n^2)
// 空间复杂度: O(n)
func minCost(costs [][]int) int {
	var (
		dfs func(i, exclude, sum int)
		ans = math.MaxInt32
	)
	n, m := len(costs), len(costs[0])
	dfs = func(i, exclude, sum int) {
		if i == n {
			ans = min(ans, sum)
			return
		}
		for j := 0; j < m; j++ {
			if j == exclude {
				continue
			}
			dfs(i+1, j, sum+costs[i][j])
		}
	}
	dfs(0, -1, 0)
	return ans
}

// minCost_2 动态规划
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func minCost_2(costs [][]int) int {
	if len(costs) == 0 {
		return 0
	}

	n, m := len(costs), len(costs[0])
	dp := costs[0]

	for i := 1; i < n; i++ {
		dpNew := make([]int, m)
		for j := 0; j < m; j++ {
			dpNew[j] = math.MaxInt32
			for k := 0; k < m; k++ {
				if k == j {
					continue
				}
				dpNew[j] = min(dpNew[j], costs[i][j]+dp[k])
			}
		}
		dp = dpNew
	}

	ans := math.MaxInt32
	for i := range dp {
		ans = min(ans, dp[i])
	}
	return ans
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
