package leetcode_0159_longest_substring_with_at_most_two_distinct_characters

// 0159.至多包含两个不同字符的最长子串
// https://leetcode-cn.com/problems/longest-substring-with-at-most-two-distinct-characters/

// lengthOfLongestSubstringTwoDistinct 滑动窗口+双指针
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func lengthOfLongestSubstringTwoDistinct(s string) int {
	n := len(s)
	if n < 3 {
		return n
	}

	m := map[byte]int{} // 记录字符出现次数
	ans := 0
	for left, right := 0, 0; right < n; right++ {
		m[s[right]]++
		for len(m) > 2 {
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
