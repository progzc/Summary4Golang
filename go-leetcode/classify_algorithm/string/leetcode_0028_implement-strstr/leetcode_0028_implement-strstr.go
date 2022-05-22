package leetcode_0028_implement_strstr

import "strings"

// 0028.实现strStr()
// https://leetcode.cn/problems/implement-strstr/

// strStr 直接使用API
// 时间复杂度: O()
// 空间复杂度: O()
func strStr(haystack string, needle string) int {
	if len(needle) == 0 {
		return 0
	}

	if len(haystack) == 0 {
		return -1
	}

	return strings.Index(haystack, needle)
}

// strStr_2 暴力法
// 时间复杂度: O((m-n)*n)
// 空间复杂度: O(m-n)
func strStr_2(haystack string, needle string) int {
	m, n := len(haystack), len(needle)
	if n == 0 {
		return 0
	}
	if m == 0 {
		return -1
	}
	for i := 0; i < m-n+1; i++ {
		find := true
		for k := range needle {
			if i+k >= m || haystack[i+k] != needle[k] {
				find = false
				break
			}
		}
		if find {
			return i
		}
	}
	return -1
}
