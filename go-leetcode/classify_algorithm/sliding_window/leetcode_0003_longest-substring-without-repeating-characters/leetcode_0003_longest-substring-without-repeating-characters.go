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

	m := map[byte]int{}
	ans := 0
	for l, r := 0, 0; r < n; r++ {
		m[s[r]]++
		for r-l+1 > len(m) {
			m[s[l]]--
			if m[s[l]] == 0 {
				delete(m, s[l])
			}
			l++
		}
		ans = max(ans, r-l+1)
	}
	return ans
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
