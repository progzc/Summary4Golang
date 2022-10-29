package leetcode_0290_word_pattern

import "strings"

// 0290. 单词规律
// https://leetcode.cn/problems/word-pattern/

// wordPattern 哈希表
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 注意: 双向匹配
// 例1:
//	输入: pattern="abba", s="dog dog dog dog"
//	输出: false
func wordPattern(pattern string, s string) bool {
	ss := strings.Split(s, " ")
	pLen, sLen := len(pattern), len(ss)
	if pLen != sLen {
		return false
	}

	m := make(map[byte]string)
	set := make(map[string]bool)
	for i := range pattern {
		if v, ok := m[pattern[i]]; !ok {
			if !set[ss[i]] {
				m[pattern[i]] = ss[i]
				set[ss[i]] = true
			} else {
				return false
			}
		} else {
			if v != ss[i] {
				return false
			}
		}
	}

	return true
}
