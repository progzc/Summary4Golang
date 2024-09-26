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

	m := map[byte]int{} // 记录每个字符出现的次数
	ans := 0
	for l, r := 0, 0; r < n; r++ {
		m[s[r]]++
		for r-l+1 > len(m) {
			m[s[l]]--
			if m[s[l]] == 0 {
				delete(m, s[l])
			}
			l++ // 注意 l++ 的位置
		}
		ans = max(ans, r-l+1)
	}
	return ans
}

// lengthOfLongestSubstring_2 滑动窗口
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func lengthOfLongestSubstring_2(s string) int {
	n := len(s)
	if n < 2 {
		return n
	}

	m := map[byte]int{} // 记录每个字符出现的位置
	ans := 0
	for l, r := 0, 0; r < n; r++ {
		if pos, ok := m[s[r]]; ok {
			// 通过 s = "abba" 示例来理解这里为何要加 max 条件
			// 本质是保证 l 永远不会倒退
			l = max(l, pos+1)
		}
		m[s[r]] = r
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
