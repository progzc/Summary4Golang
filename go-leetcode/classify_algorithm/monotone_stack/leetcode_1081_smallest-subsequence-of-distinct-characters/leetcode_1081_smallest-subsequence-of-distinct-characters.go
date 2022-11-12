package leetcode_1081_smallest_subsequence_of_distinct_characters

// 1081. 不同字符的最小子序列
// https://leetcode.cn/problems/smallest-subsequence-of-distinct-characters/

// 题目同：
// 0316. 去除重复字母
// https://leetcode.cn/problems/remove-duplicate-letters/

// smallestSubsequence 单调递增(非严格)栈
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路:
//	i)先遍历字符串，记录每个元素出现次数
//	ii)接着元素挨个入栈, 维持栈的核心要素是：保持栈内所有元素都小于即将进栈的元素，即字典序；
//	   如果不满足，则需要把栈内元素出栈，直到满足要求为止。
func smallestSubsequence(s string) string {
	n := len(s)
	if n <= 1 {
		return s
	}

	var (
		// count 统计每个字符的出现次数
		count = [26]int{}
		// exist 是否在栈中
		exist = [26]bool{}
		// stack 单调栈
		stack []byte
	)

	for i := 0; i < n; i++ {
		count[s[i]-'a']++
	}

	for i := 0; i < n; i++ {
		if exist[s[i]-'a'] {
			count[s[i]-'a']--
			continue
		}
		// 注意: for len(stack) > 0 && stack[len(stack)-1] > s[i] && count[s[i]-'a'] > 0 错误
		for len(stack) > 0 && stack[len(stack)-1] > s[i] && count[stack[len(stack)-1]-'a'] > 0 {
			exist[stack[len(stack)-1]-'a'] = false
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, s[i])
		exist[s[i]-'a'] = true
		count[s[i]-'a']--
	}
	return string(stack)
}
