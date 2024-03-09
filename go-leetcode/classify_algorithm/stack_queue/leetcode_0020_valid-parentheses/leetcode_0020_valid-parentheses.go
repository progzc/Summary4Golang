package leetcode_0020_valid_parentheses

import "strings"

// 0020.有效的括号
// https://leetcode-cn.com/problems/valid-parentheses/

// isValid 栈
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路：经典的栈解法
func isValid_1(s string) bool {
	n := len(s)
	if n%2 == 1 {
		return false
	}
	pairMap := map[byte]byte{
		')': '(',
		']': '[',
		'}': '{',
	}
	stack := []byte{}
	for i := 0; i < n; i++ {
		if _, ok := pairMap[s[i]]; ok {
			if len(stack) == 0 || stack[len(stack)-1] != pairMap[s[i]] {
				return false
			}
			stack = stack[:len(stack)-1]
		} else {
			stack = append(stack, s[i])
		}
	}
	return len(stack) == 0
}

// isValid_3 栈
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路：经典的栈解法
func isValid_3(s string) bool {
	n := len(s)
	if n%2 != 0 {
		return false
	}
	var stack []byte
	for i := 0; i < n; i++ {
		if s[i] == '(' || s[i] == '[' || s[i] == '{' {
			stack = append(stack, s[i])
		} else {
			if len(stack) == 0 {
				return false
			}
			pre := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if (s[i] == ')' && pre != '(') || (s[i] == ']' && pre != '[') || (s[i] == '}' && pre != '{') {
				return false
			}
		}
	}
	return len(stack) == 0
}

// isValid_2 替代
// 时间复杂度: O(n^2)
// 空间复杂度: O(n)
// 思路: 消消乐
func isValid_2(s string) bool {
	sLen := 0
	for sLen != len(s) {
		sLen = len(s)
		s = strings.ReplaceAll(s, "()", "")
		s = strings.ReplaceAll(s, "{}", "")
		s = strings.ReplaceAll(s, "[]", "")
	}
	return len(s) == 0
}
