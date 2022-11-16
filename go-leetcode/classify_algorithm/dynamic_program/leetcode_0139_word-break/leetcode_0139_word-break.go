package leetcode_0139_word_break

import "sort"

// 0139.单词拆分
// https://leetcode-cn.com/problems/word-break/

// wordBreak 动态规划
// 时间复杂度：O(n^2)
// 空间复杂度：O(n)
// 	定义dp[i]表示字符串s的前i个字符组成的字符串s[0..i-1]是否能被空格拆分成若干个字典中出现的单词，则
//		dp[i]=dp[j]&&check(s[j..i-1])
// 注意：不要求字典中出现的单词全部都使用，并且字典中的单词可以重复使用。
func wordBreak(s string, wordDict []string) bool {
	set := make(map[string]bool)
	for _, w := range wordDict {
		set[w] = true
	}
	dp := make([]bool, len(s)+1)
	dp[0] = true
	for i := 0; i <= len(s); i++ {
		for j := 0; j < i; j++ {
			if dp[j] && set[s[j:i]] {
				dp[i] = true
				break
			}
		}
	}
	return dp[len(s)]
}

// wordBreak_2 dfs (超时)
func wordBreak_2(s string, wordDict []string) bool {
	sort.Slice(wordDict, func(i, j int) bool {
		return len(wordDict[i]) < len(wordDict[j])
	})
	var (
		dfs func(begin int, str string) bool
		n   = len(wordDict)
	)
	dfs = func(begin int, str string) bool {
		if begin < n {
			if str == s {
				return true
			}
		}

		for i := begin; i < n; i++ {
			if len(wordDict[i])+len(str) > len(s) {
				break
			}
			if dfs(begin, str+wordDict[i]) {
				return true
			}
		}
		return false
	}
	return dfs(0, "")
}
