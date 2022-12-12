package leetcode_1003_check_if_word_is_valid_after_substitutions

import "strings"

// 1003. 检查替换后的词是否有效
// https://leetcode.cn/problems/check-if-word-is-valid-after-substitutions/

// isValid
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func isValid(s string) bool {
	n := len(s)
	if n == 0 || n%3 != 0 {
		return false
	}

	for len(s) > 0 {
		newS := strings.ReplaceAll(s, "abc", "")
		if len(newS) == len(s) {
			break
		}
		s = newS
	}
	return len(s) == 0
}

// isValid 栈
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路:
//	a.碰到c就出a、b 不是a、b直接报错。
//	b.碰到a,b就入栈。
//	c.循环完了栈长度大于0,则返回false
func isValid_2(s string) bool {
	n := len(s)
	if n == 0 || n%3 != 0 {
		return false
	}
	var stack []byte
	for _, ch := range s {
		if ch == 'c' {
			if len(stack) >= 2 && stack[len(stack)-1] == 'b' && stack[len(stack)-2] == 'a' {
				stack = stack[:len(stack)-2]
			} else {
				return false
			}
		} else {
			stack = append(stack, byte(ch))
		}
	}
	return len(stack) == 0
}
