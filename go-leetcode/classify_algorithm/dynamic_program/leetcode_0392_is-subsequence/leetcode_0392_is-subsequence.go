package leetcode_0392_is_subsequence

// 392. 判断子序列
// https://leetcode.cn/problems/is-subsequence/

// 进阶:
// 如果有大量输入的 S，称作 S1, S2, ... , Sk 其中 k >= 10亿，你需要依次检查它们是否为 T 的子序列。在这种情况下，你会怎样改变代码？

// isSubsequence 双指针
// 时间复杂度: O(m+n)
// 空间复杂度: O(1)
func isSubsequence(s string, t string) bool {
	sLen, tLen := len(s), len(t)
	if sLen == 0 {
		return true
	}
	if tLen == 0 || sLen > tLen {
		return false
	}

	p1, p2 := 0, 0
	for p1 < sLen && p2 < tLen {
		for p2 < tLen && t[p2] != s[p1] {
			p2++
		}
		if p2 == tLen {
			return false
		}
		p1++
		p2++
	}
	return p1 == sLen
}

// isSubsequence_2 动态规划
// 时间复杂度: O(mn)
// 空间复杂度: O(mn)
// 思路:
//	状态: dp[i][j]表示s[0...i]是否是t[0...j]的子序列
//	转移方程:
//	  若s[i]==t[j], dp[i][j] = dp[i-1][j-1]
//	  若s[i]!=t[j], dp[i][j] = dp[i][j-1]
func isSubsequence_2(s string, t string) bool {
	sLen, tLen := len(s), len(t)
	if sLen == 0 {
		return true
	}
	if tLen == 0 || sLen > tLen {
		return false
	}

	dp := make([][]bool, sLen)
	for i := 0; i < sLen; i++ {
		dp[i] = make([]bool, tLen)
	}
	dp[0][0] = s[0] == t[0]
	for j := 1; j < tLen; j++ {
		dp[0][j] = dp[0][j-1] || s[0] == t[j]
	}

	for i := 1; i < sLen; i++ {
		// j从i开始，这是因为当i>j时，肯定有dp[i][j]=false
		for j := i; j < tLen; j++ {
			if s[i] == t[j] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = dp[i][j-1]
			}
		}
	}

	return dp[sLen-1][tLen-1]
}

// isSubsequence_3 动态规划（解决进阶问题）
// 时间复杂度: O(mn)
// 空间复杂度: O(mn)
// 思路: 可以预处理出对于 t 的每一个位置，从该位置开始往后每一个字符第一次出现的位置
//	状态: f[i][j]表示字符串 t 中从位置 i 开始往后字符 j 第一次出现的位置。
//	转移方程:
//	  若t[i]==j, f[i][j] = i
//	  若t[i]!=j, f[i][j] = f[i+1][j]
func isSubsequence_3(s string, t string) bool {
	n, m := len(s), len(t)
	f := make([][26]int, m+1)
	for i := 0; i < 26; i++ {
		f[m][i] = m
	}
	for i := m - 1; i >= 0; i-- {
		for j := 0; j < 26; j++ {
			if t[i] == byte(j+'a') {
				f[i][j] = i
			} else {
				f[i][j] = f[i+1][j]
			}
		}
	}
	add := 0
	for i := 0; i < n; i++ {
		if f[add][int(s[i]-'a')] == m {
			return false
		}
		add = f[add][int(s[i]-'a')] + 1
	}
	return true
}
