package leetcode_0003_longest_substring_without_repeating_characters

// 0003.无重复字符的最长子串
// https://leetcode-cn.com/problems/longest-substring-without-repeating-characters/

// lengthOfLongestSubstring 滑动窗口
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func lengthOfLongestSubstring(s string) int {
	n := len(s)
	if n < 2 {
		return n
	}

	rk := 0
	m := map[byte]int{}
	ans := 0
	for i := 0; i < n; i++ {
		if i != 0 {
			delete(m, s[i-1])
		}
		for rk < n && m[s[rk]] == 0 {
			m[s[rk]]++
			rk++
		}
		ans = max(ans, rk-i)
	}
	return ans
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
