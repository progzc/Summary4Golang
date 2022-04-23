package leetcode_1143_longest_common_subsequence

// 1143.最长公共子序列
// https://leetcode-cn.com/problems/longest-common-subsequence/

// longestCommonSubsequence 二维动态规划
// 时间复杂度: O(m*n)
// 空间复杂度: O(m*n)
// 思路：典型的动态规划
func longestCommonSubsequence(text1 string, text2 string) int {
	m, n := len(text1), len(text2)
	// dp[i][j]: 表示text1[:i]与text2[:j]的最长公共子序列的长度
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	// 初始化: dp[0][j]=0, dp[i][0]=0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if text1[i] == text2[j] {
				dp[i+1][j+1] = dp[i][j] + 1
			} else {
				dp[i+1][j+1] = max(dp[i][j+1], dp[i+1][j])
			}
		}
	}
	return dp[m][n]
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
