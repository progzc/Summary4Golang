package leetcode_0014_longest_common_prefix

import (
	"bytes"
	"strings"
)

// 0014.最长公共前缀
// https://leetcode.cn/problems/longest-common-prefix/

// longestCommonPrefix
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func longestCommonPrefix(strs []string) string {
	i, ans := 0, bytes.Buffer{}
	for {
		for j := range strs {
			if i >= len(strs[j]) || strs[0][i] != strs[j][i] {
				return ans.String()
			}
		}
		ans.WriteByte(strs[0][i])
		i++
	}
}

// longestCommonPrefix_2 零拷贝优化
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func longestCommonPrefix_2(strs []string) string {
	i, ans := 0, strings.Builder{}
	for {
		for j := range strs {
			if i >= len(strs[j]) || strs[0][i] != strs[j][i] {
				return ans.String()
			}
		}
		ans.WriteByte(strs[0][i])
		i++
	}
}

// longestCommonPrefix_3 空间进一步优化
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func longestCommonPrefix_3(strs []string) string {
	i := 0
	for {
		for j := range strs {
			if i >= len(strs[j]) || strs[0][i] != strs[j][i] {
				return strs[0][:i]
			}
		}
		i++
	}
}
