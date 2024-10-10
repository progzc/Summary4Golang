package leetcode_1143_longest_common_subsequence

// 1143.最长公共子序列🌟
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

// longestCommonSubsequence_2 二维动态规划
// 时间复杂度: O(m*n)
// 空间复杂度: O(m*n)
// 思路：典型的动态规划
func longestCommonSubsequence_2(text1 string, text2 string) int {
	m, n := len(text1), len(text2)
	// dp[i][j]: 表示text1[:i]与text2[:j]的最长公共子序列的长度
	dp := make([][]int, m)
	for i := range dp {
		dp[i] = make([]int, n)
	}
	// 初始化:
	for i := 0; i < m; i++ {
		if i == 0 {
			if text1[i] == text2[0] {
				dp[i][0] = 1
			}
		} else {
			if text1[i] == text2[0] || dp[i-1][0] > 0 {
				dp[i][0] = 1
			}
		}
	}

	for j := 0; j < n; j++ {
		if j == 0 {
			if text1[0] == text2[j] {
				dp[0][j] = 1
			}
		} else {
			if text1[0] == text2[j] || dp[0][j-1] > 0 {
				dp[0][j] = 1
			}
		}
	}

	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			if text1[i] == text2[j] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}
	return dp[m-1][n-1]
}

// longestCommonSubsequence_3 二维动态规划
// 时间复杂度: O(m*n)
// 空间复杂度: O(m*n)
// 思路：典型的动态规划
func longestCommonSubsequence_3(text1 string, text2 string) int {
	m, n := len(text1), len(text2)
	dp := make([][]int, m)
	for i := 0; i < m; i++ {
		dp[i] = make([]int, n)
		if i == 0 {
			if text1[i] == text2[0] {
				dp[i][0] = 1
			}
		} else {
			if text1[i] == text2[0] {
				dp[i][0] = 1
			} else {
				dp[i][0] = dp[i-1][0]
			}
		}
	}

	for j := 0; j < n; j++ {
		if j == 0 {
			if text1[0] == text2[j] {
				dp[0][j] = 1
			}
		} else {
			if text1[0] == text2[j] {
				dp[0][j] = 1
			} else {
				dp[0][j] = dp[0][j-1]
			}
		}
	}

	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			if text1[i] == text2[j] {
				dp[i][j] = max(dp[i-1][j-1]+1, max(dp[i-1][j], dp[i][j-1]))
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}
	return dp[m-1][n-1]
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
