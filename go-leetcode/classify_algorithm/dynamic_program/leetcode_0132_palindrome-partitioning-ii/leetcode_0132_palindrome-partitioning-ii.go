package leetcode_0132_palindrome_partitioning_ii

import "math"

// 132. 分割回文串 II
// https://leetcode.cn/problems/palindrome-partitioning-ii/

// minCut dfs (超时)
// 时间复杂度: O(n*2^n)
// 空间复杂度: O(n^2)
func minCut(s string) int {
	var isPalindrome func(s string) bool
	isPalindrome = func(s string) bool {
		sLen := len(s)
		for i := 0; i < sLen/2; i++ {
			if s[i] != s[sLen-i-1] {
				return false
			}
		}
		return true
	}

	n := len(s)
	if n <= 1 || isPalindrome(s) {
		return 0
	}
	minSplit := math.MaxInt32
	for i := 1; i <= n-1; i++ {
		if isPalindrome(s[:i]) {
			minSplit = min(minSplit, minCut(s[i:]))
		}
	}
	return minSplit + 1
}

// minCut_2 dfs + 动态规划预处理 (仍然超时)
// 时间复杂度: O(n*2^n)
// 空间复杂度: O(n^2)
// 思路: 判断字符串是否为回文字符串可以使用动态规划进行预处理，可以加快速度。
//	状态: dp[i][j]表示s[i...j](左闭右闭区间)是否为回文串。
//	转移方程:
//		当j-i<0时, dp[i][j] = false
//		当j-i==0时,dp[i][j] = true
//		当j-i==1时,dp[i][j] = s[i] == s[j]
//		当j-i>=2时,dp[i][j] = s[i] == s[j] && dp[i+1][j-1]
func minCut_2(s string) int {
	n := len(s)
	dp := make([][]bool, n)
	for i := 0; i < n; i++ {
		dp[i] = make([]bool, n)
	}
	for i := n - 1; i >= 0; i-- {
		for j := i; j < n; j++ {
			if j-i == 0 {
				dp[i][j] = true
			} else if j-i == 1 {
				dp[i][j] = s[i] == s[j]
			} else {
				dp[i][j] = s[i] == s[j] && dp[i+1][j-1]
			}
		}
	}

	// dfs 求解s[i...j](左闭右闭)的最少分割次数
	var dfs func(i, j int) int
	dfs = func(i, j int) int {
		if dp[i][j] {
			return 0
		}
		minSplit := math.MaxInt32
		for end := i; end <= j-1; end++ {
			if dp[i][end] {
				minSplit = min(minSplit, dfs(end+1, j))
			}
		}
		return minSplit + 1
	}
	return dfs(0, n-1)
}

// minCut_3 dfs + 动态规划预处理 (优化)
// 时间复杂度: O(n*2^n)
// 空间复杂度: O(n^2)
// 思路: 判断字符串是否为回文字符串可以使用动态规划进行预处理，可以加快速度。
//	状态一: dp[i][j]表示s[i...j](左闭右闭区间)是否为回文串。
//	转移方程:
//		当j-i<0时, dp[i][j] = false
//		当j-i==0时,dp[i][j] = true
//		当j-i==1时,dp[i][j] = s[i] == s[j]
//		当j-i>=2时,dp[i][j] = s[i] == s[j] && dp[i+1][j-1]
//
//  状态二: f[i]表示字符串的前缀 s[0..i] 的最少分割次数
//	转移方程:
//		f[i] = min{f[0], f[1], ...,f[j]}+1,  其中0<=j<i，且s[j+1...i]是一个回文字符串
func minCut_3(s string) int {
	n := len(s)
	dp := make([][]bool, n)
	for i := 0; i < n; i++ {
		dp[i] = make([]bool, n)
	}
	for i := n - 1; i >= 0; i-- {
		for j := i; j < n; j++ {
			if j-i == 0 {
				dp[i][j] = true
			} else if j-i == 1 {
				dp[i][j] = s[i] == s[j]
			} else {
				dp[i][j] = s[i] == s[j] && dp[i+1][j-1]
			}
		}
	}

	f := make([]int, n)
	for i := 0; i < n; i++ {
		if dp[0][i] {
			f[i] = 0
		} else {
			f[i] = math.MaxInt32
			for j := 0; j < i; j++ {
				if dp[j+1][i] && f[j]+1 < f[i] {
					f[i] = f[j] + 1
				}
			}
		}
	}
	return f[n-1]
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
