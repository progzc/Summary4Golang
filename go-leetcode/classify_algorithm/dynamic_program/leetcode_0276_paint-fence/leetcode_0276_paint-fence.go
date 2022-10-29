package leetcode_0276_paint_fence

// 0276. 栅栏涂色
// https://leetcode.cn/problems/paint-fence/

// numWays dfs（超时）
// 时间复杂度: O(n^k)
// 空间复杂度: O(n)
func numWays(n int, k int) int {
	if n == 1 {
		return k
	}

	var dfs func(i int, plans []int)
	count := 0
	dfs = func(i int, plans []int) {
		if i == n {
			count++
			return
		}

		for j := 1; j <= k; j++ {
			if len(plans) >= 2 && plans[len(plans)-1] == j && plans[len(plans)-2] == j {
				continue
			}
			plans = append(plans, j)
			dfs(i+1, plans)
			plans = plans[:len(plans)-1]
		}
	}
	dfs(0, nil)
	return count
}

// numWays_2 动态规划（二维dp）
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路:
// 	状态: 由于每个栅栏都有两个状态： 与上一个颜色相同、与上一个颜色不同。 那么我们可以开辟一个 dp[n][2] 大小的数组来表示：
//		 dp[i][0]: 状态 0 表示与上一个颜色相同，dp[i][1]: 状态 1 表示与上一个颜色不同
//	转移方程:
//		dp[i][0] = dp[i - 1][1]
//		dp[i][1] = (dp[i - 1][0] + dp[i - 1][1]) * (k - 1)
//	最终结果为:
//		dp[n-1][0] + dp[n-1][1]
func numWays_2(n int, k int) int {
	dp := make([][2]int, n)
	dp[0][0] = 0
	dp[0][1] = k
	for i := 1; i < n; i++ {
		dp[i][0] = dp[i-1][1]
		dp[i][1] = (dp[i-1][0] + dp[i-1][1]) * (k - 1)
	}
	return dp[n-1][0] + dp[n-1][1]
}

// numWays_3 动态规划（一维dp）
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路:
// 	状态: 由于每个栅栏都有两个状态： 与上一个颜色相同、与上一个颜色不同。 那么我们可以开辟一个 dp[n][2] 大小的数组来表示：
//		 dp[i][0]: 状态 0 表示与上一个颜色相同，dp[i][1]: 状态 1 表示与上一个颜色不同
//	转移方程:
//		dp[i][0] = dp[i - 1][1]
//		dp[i][1] = (dp[i - 1][0] + dp[i - 1][1]) * (k - 1)
//	最终结果为:
//		dp[n-1][0] + dp[n-1][1]
func numWays_3(n int, k int) int {
	dp0, dp1 := 0, k
	for i := 1; i < n; i++ {
		dp0, dp1 = dp1, (dp0+dp1)*(k-1)
	}
	return dp0 + dp1
}
