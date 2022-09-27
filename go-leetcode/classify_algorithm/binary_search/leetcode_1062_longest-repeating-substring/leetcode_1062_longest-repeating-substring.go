package leetcode_1062_longest_repeating_substring

import (
	"strings"
)

// 1062. 最长重复子串
// https://leetcode.cn/problems/longest-repeating-substring/

// TODO 二分查找的解法

// longestRepeatingSubstring 动态规划
// 时间复杂度: O(n^2)
// 空间复杂度: O(n^2)
// 思路：
//	状态定义为dp[i][j]是两个分别以i和j结尾的相同子串的最大长度，其中i永远小于j，所有状态的值均初始化为0。
//	状态转移时，如果s[i]和s[j]不同就不必管，因为以i结尾和以j结尾不会是相同子串。
//	如果s[i]和s[j]相同，那么dp[i][j]就等于dp[i-1][j-1]+1，这点应该是很显然的，就是给i-1和j-1结尾的重复子串两边各加了一个相同字符。
//	注意此时如果i=0，那么dp[i][j]就是1
func longestRepeatingSubstring(s string) int {
	n := len(s)
	dp := make([][]int, n)
	for i := range dp {
		dp[i] = make([]int, n)
	}

	ans := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if s[i] == s[j] {
				if i == 0 {
					dp[i][j] = 1
				} else {
					dp[i][j] = dp[i-1][j-1] + 1
				}
				ans = max(ans, dp[i][j])
			}
		}
	}

	return ans
}

// longestRepeatingSubstring_2 滑动窗口
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func longestRepeatingSubstring_2(s string) int {
	begin, end := 0, 1
	n, ans := len(s), 0
	for end < n {
		subStr := s[begin:end]
		if strings.Index(s[begin+1:], subStr) != -1 {
			ans = max(ans, end-begin)
			end++
		} else {
			begin++
			if begin == end {
				end++
			}
		}
	}
	return ans
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
