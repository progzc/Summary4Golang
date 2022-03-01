package leetcode_0242_valid_anagram

import "sort"

// 242.有效的字母异位词
// https://leetcode-cn.com/problems/valid-anagram/

// isAnagram_1 哈希表（借助数组）
// 时间复杂度: O(n)
// 空间复杂度: O(s)
func isAnagram_1(s string, t string) bool {
	var c1, c2 [26]int
	for _, ch := range s {
		c1[ch-'a']++
	}
	for _, ch := range t {
		c2[ch-'a']++
	}
	// golang中两个数组相同的条件是：
	// 1. 类型相同（数组长度是类型的一部分）
	// 2. 数组元素相同
	return c1 == c2
}

// isAnagram_2 哈希表（借助map）
// 时间复杂度: O(n)
// 空间复杂度: O(s)
func isAnagram_2(s string, t string) bool {
	if len(s) != len(t) {
		return false
	}
	cnt := map[rune]int{}
	for _, ch := range s {
		cnt[ch]++
	}
	for _, ch := range t {
		cnt[ch]--
		if cnt[ch] < 0 {
			return false
		}
	}
	return true
}

// isAnagram_3 排序
// 时间复杂度: O(nlog(n))
// 空间复杂度: O(log(n))
func isAnagram_3(s string, t string) bool {
	s1, s2 := []byte(s), []byte(t)
	sort.Slice(s1, func(i, j int) bool {
		return s1[i] < s1[j]
	})
	sort.Slice(s2, func(i, j int) bool {
		return s2[i] < s2[j]
	})
	return string(s1) == string(s2)
}
