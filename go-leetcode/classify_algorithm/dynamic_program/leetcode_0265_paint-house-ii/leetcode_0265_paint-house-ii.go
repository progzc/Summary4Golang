package leetcode_0265_paint_house_ii

import "math"

// 265. 粉刷房子 II
// https://leetcode.cn/problems/paint-house-ii/

// minCostII 动态规划
// 时间复杂度: O(nk^2)
// 空间复杂度: O(nk)
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

// minCostII_2 动态规划
// 时间复杂度: O(nk)
// 空间复杂度: O(nk)
// 思路: 类似于[198. 打家劫舍]
//	状态: dp[i][j]表示粉刷第 0 号房子到第 i 号房子且第 i 号房子被粉刷成第 j 种颜色时的最小花费成本。
//	转移方程: dp[i][j] = min(dp[i-1][0],...,dp[i-1][j-1],dp[i-1][j+1]...,dp[i-1][n-1])+cost[i][j]
// 注意:
//	 此题代码可同样用于[256.粉刷房子]
// 	 要保证时间复杂度为O(nk), 可以记录当前房子的一种颜色的花费的最小值和第二小值,为下一个房子做准备。
func minCostII_2(costs [][]int) int {
	if len(costs) == 0 {
		return 0
	}

	n, m := len(costs), len(costs[0])
	dp := make([][]int, n)
	for i := 0; i < n; i++ {
		dp[i] = make([]int, m)
	}
	// c1表示上一个房子粉刷为某种颜色的最低花费, c2表示上一个房子粉刷为某种颜色的第二低花费
	c1, c2 := math.MaxInt32, math.MaxInt32
	for j := 0; j < m; j++ {
		dp[0][j] = costs[0][j]
		if dp[0][j] < c1 {
			c2, c1 = c1, dp[0][j]
		} else if dp[0][j] < c2 {
			c2 = dp[0][j]
		}
	}

	for i := 1; i < n; i++ {
		tmp1, tmp2 := math.MaxInt32, math.MaxInt32
		for j := 0; j < m; j++ {
			// 如果当前颜色j与上一个房子的花费最小颜色花费不一样，就可以直接取上一个房子颜色的最低花费，与当前花费形成最优解；
			// 如果一样，就取上一个房子的第二低花费来形成当前颜色的最优解。
			if dp[i-1][j] != c1 {
				dp[i][j] = costs[i][j] + c1
			} else {
				dp[i][j] = costs[i][j] + c2
			}
			// 在得到当前房子的最优解时候同时为下一个房子求出最低花费和第二低花费；
			// 就是维护一个最小值，和一个第二小的值。
			if dp[i][j] < tmp1 {
				tmp2, tmp1 = tmp1, dp[i][j]
			} else if dp[i][j] < tmp2 {
				tmp2 = dp[i][j]
			}
		}
		// 更新c1,c2
		c1, c2 = tmp1, tmp2
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
