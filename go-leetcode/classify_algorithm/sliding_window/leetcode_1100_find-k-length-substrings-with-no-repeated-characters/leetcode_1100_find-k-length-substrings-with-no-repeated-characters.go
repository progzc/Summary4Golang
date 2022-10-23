package leetcode_1100_find_k_length_substrings_with_no_repeated_characters

// 1100. 长度为 K 的无重复字符子串
// https://leetcode.cn/problems/find-k-length-substrings-with-no-repeated-characters/

// 可以对比 https://leetcode.cn/problems/longest-repeating-character-replacement/

// numKLenSubstrNoRepeats 滑动窗口
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func numKLenSubstrNoRepeats(s string, k int) int {
	n := len(s)
	if k > n {
		return 0
	}

	m := make(map[byte]int)
	count := 0
	for i := 0; i < k; i++ {
		m[s[i]]++
	}
	if len(m) == k {
		count++
	}

	for i := k; i < n; i++ {
		m[s[i-k]]--
		if m[s[i-k]] == 0 {
			delete(m, s[i-k])
		}
		m[s[i]]++
		if len(m) == k {
			count++
		}
	}
	return count
}

// numKLenSubstrNoRepeats_2 双指针
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func numKLenSubstrNoRepeats_2(s string, k int) int {
	n := len(s)
	if k > n {
		return 0
	}

	count := 0
	m := make(map[byte]int)
	for l, r := 0, 0; r < n; r++ {
		m[s[r]]++
		// 当数量超过k 或者 出现重复时，移动左指针
		for len(m) > k || m[s[r]] > 1 {
			m[s[l]]--
			if m[s[l]] == 0 {
				delete(m, s[l])
			}
			l++
		}
		if len(m) == k {
			count++
		}
	}
	return count
}
