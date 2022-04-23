package leetcode_0516_longest_palindromic_subsequence

// 0516.最长回文子序列
// https://leetcode-cn.com/problems/longest-palindromic-subsequence/

// longestPalindromeSubseq 动态规划（区间DP）
// 时间复杂度: O(n*n)
// 空间复杂度: O(n*n)
// 状态: dp[i][j]表示字符串s的下标范围[i,j]内的最长回文子序列的长度
// 边界条件：
//	a.当0<=i<=j<n时，才会有dp[i][j]>0，否则dp[i][j]=0
//	b.对于任意的0<=i<n，都有dp[i][i]=1
// 转移方程：
//	当i<j，s[i]=s[j]时，dp[i][j]=dp[i+1][j-1]+2
//	当i<j，s[i]!=s[j]时，dp[i][j]=max(dp[i+1][j],dp[i][j-1])
// 注意事项：由于状态转移方程都是从长度较短的子序列向长度较长的子序列转移，因此需要注意动态规划的循环顺序。
func longestPalindromeSubseq(s string) int {
	n := len(s)
	dp := make([][]int, n)
	for i := range dp {
		dp[i] = make([]int, n)
	}
	// i必须逆序
	for i := n - 1; i >= 0; i-- {
		dp[i][i] = 1
		for j := i + 1; j < n; j++ {
			if s[i] == s[j] {
				dp[i][j] = dp[i+1][j-1] + 2
			} else {
				dp[i][j] = max(dp[i+1][j], dp[i][j-1])
			}
		}
	}
	return dp[0][n-1]
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
