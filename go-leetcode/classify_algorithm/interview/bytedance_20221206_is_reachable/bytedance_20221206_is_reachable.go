package bytedance_20221206_is_reachable

// isReachable 寻路算法
// 时间复杂度: O(n^2)
// 空间复杂度: O(n^2)
func isReachable(nums [][]int) bool {
	if len(nums) == 0 || len(nums[0]) == 0 {
		return false
	}
	if nums[0][0] == 0 {
		return false
	}

	m, n := len(nums), len(nums[0])
	dp := make([][]bool, m)
	for i := 0; i < m; i++ {
		dp[i] = make([]bool, n)
	}
	dp[0][0] = true

	for i := 1; i < m; i++ {
		dp[i][0] = nums[i][0] == 1 && dp[i-1][0]
	}
	for j := 1; j < n; j++ {
		dp[0][j] = nums[0][j] == 1 && dp[0][j-1]
	}

	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			dp[i][j] = nums[i][j] == 1 && (dp[i-1][j] || dp[i][j-1])
		}
	}
	return dp[m-1][n-1]
}
