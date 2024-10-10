package leetcode_0062_unique_paths

// 0062. 不同路径
// https://leetcode.cn/problems/unique-paths

// uniquePaths 动态规划
// 时间复杂度: O(mn)
// 空间复杂度: O(mn)
func uniquePaths(m int, n int) int {
	dp := make([][]int, m)
	for i := 0; i < m; i++ {
		dp[i] = make([]int, n)
		dp[i][0] = 1
	}
	for j := 0; j < n; j++ {
		dp[0][j] = 1
	}
	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			dp[i][j] = dp[i][j-1] + dp[i-1][j]
		}
	}
	return dp[m-1][n-1]
}
