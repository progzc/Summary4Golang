package leetcode_0028_implement_strstr

import "strings"

// 0028.实现strStr()
// https://leetcode.cn/problems/implement-strstr/

// KMP算法：https://zhuanlan.zhihu.com/p/83334559

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
// 时间复杂度: O(m*n)
// 空间复杂度: O(1)
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

// strStr_3 KMP
// 时间复杂度: O(m+n)
// 空间复杂度: O(n)
// 思路：空间换时间
func strStr_3(haystack string, needle string) int {
	m, n := len(haystack), len(needle)
	if n == 0 {
		return 0
	}
	if m == 0 {
		return -1
	}

	// 使用KMP算法进行搜索
	// 步骤1：生成next数组
	next := make([]int, n)
	for i, j := 1, 0; i < n; i++ {
		for j > 0 && needle[i] != needle[j] {
			j = next[j-1]
		}
		if needle[i] == needle[j] {
			j++
		}
		next[i] = j
	}
	// 步骤2：利用next数组加速搜索
	for i, j := 0, 0; i < m; i++ {
		for j > 0 && haystack[i] != needle[j] {
			j = next[j-1]
		}
		if haystack[i] == needle[j] {
			j++
		}
		if j == n {
			return i - n + 1
		}
	}
	return -1
}
