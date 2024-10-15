package hehe_20241015_longest_palindromic_substring

import (
	"fmt"
	"testing"
)

//最长回文子串

//给你一个字符串 s，找到 s 中最长的
//回文子串
//示例 1：
//输入：s = "babad"
//输出："bab"
//示例 2：
//输入：s = "cbbd"
//输出："bb"

func TestPalindrome(t *testing.T) {
	fmt.Println(palindrome("babad"))
	fmt.Println(palindrome("cbbd"))
}

func palindrome(s string) string {
	n := len(s)
	if n <= 1 {
		return s
	}
	dp := make([][]bool, n)
	for i := 0; i < n; i++ {
		dp[i] = make([]bool, n)
	}

	//dp[i][j] == dp[i+1][j-1] && s[i]==s[j] or s[i]==s[j]
	maxLen := 0
	start, end := 0, 0
	for i := n - 1; i >= 0; i-- {
		j := i
		for ; j < n; j++ {
			if j-i > 1 {
				dp[i][j] = dp[i+1][j-1] && s[i] == s[j]
			} else {
				dp[i][j] = s[i] == s[j]
			}
			if dp[i][j] && j-i+1 > maxLen {
				maxLen = j - i + 1
				start = i
				end = j
			}
		}
	}
	return s[start : end+1]
}
