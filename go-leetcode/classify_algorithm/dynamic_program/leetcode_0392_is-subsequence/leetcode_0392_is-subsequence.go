package leetcode_0392_is_subsequence

// 392. 判断子序列
// https://leetcode.cn/problems/is-subsequence/

// 进阶:
// 如果有大量输入的 S，称作 S1, S2, ... , Sk 其中 k >= 10亿，你需要依次检查它们是否为 T 的子序列。在这种情况下，你会怎样改变代码？

// isSubsequence 双指针
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func isSubsequence(s string, t string) bool {
	sLen, tLen := len(s), len(t)
	if sLen == 0 {
		return true
	}
	if tLen == 0 || sLen > tLen {
		return false
	}

	p1, p2 := 0, 0
	for p1 < sLen && p2 < tLen {
		for p2 < tLen && t[p2] != s[p1] {
			p2++
		}
		if p2 == tLen {
			return false
		}
		p1++
		p2++
	}
	return p1 == sLen
}
