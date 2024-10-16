package leetcode_0340_longest_substring_with_at_most_k_distinct_characters

// 0340.至多包含 K 个不同字符的最长子串
// https://leetcode-cn.com/problems/longest-substring-with-at-most-k-distinct-characters/

// 与 https://leetcode.cn/problems/longest-repeating-character-replacement/ 类似

// lengthOfLongestSubstringKDistinct 滑动窗口+双指针
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func lengthOfLongestSubstringKDistinct(s string, k int) int {
	n := len(s)
	if n <= k {
		return n
	}

	m := map[byte]int{}
	ans := 0
	for left, right := 0, 0; right < n; right++ {
		m[s[right]]++
		for len(m) > k {
			m[s[left]]--
			if m[s[left]] == 0 {
				delete(m, s[left])
			}
			left++
		}
		ans = max(ans, right-left+1)
	}
	return ans
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
