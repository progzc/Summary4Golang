package leetcode_0020_valid_parentheses

// 0020.有效的括号
// https://leetcode-cn.com/problems/valid-parentheses/

// isValid 栈
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func isValid(s string) bool {
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
		if pairMap[s[i]] > 0 {
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
