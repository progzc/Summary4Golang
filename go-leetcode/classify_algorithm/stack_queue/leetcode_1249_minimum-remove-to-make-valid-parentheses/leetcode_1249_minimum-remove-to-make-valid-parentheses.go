package leetcode_1249_minimum_remove_to_make_valid_parentheses

// 1249.移除无效的括号
// https://leetcode-cn.com/problems/minimum-remove-to-make-valid-parentheses/

// minRemoveToMakeValid 栈
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路：
//	a.直接遍历s，遇到'('，将他的位置入栈；遇到')'有两种可能:
//		栈里有'('：就出栈;
//		没有'(': 即栈为空直接删除当前符号，此时s长度和下标需要--。
//	b.遍历完后如果栈不为空，说明存在'('没有匹配，故依次pop栈顶元素，逐个删除。
func minRemoveToMakeValid(s string) string {
	var stack []int
	for i := 0; i < len(s); i++ {
		// 直接遍历s，遇到'('，将他的位置入栈
		if s[i] == '(' {
			stack = append(stack, i)
		}
		// 遇到')'有两种可能
		if s[i] == ')' {
			if len(stack) > 0 {
				// 栈里有'('：就出栈
				stack = stack[:len(stack)-1]
			} else {
				// 没有'(': 即栈为空直接删除当前符号
				// 注意：事实上，s[i+1:]不会出现索引越界
				if i == len(s)-1 {
					s = s[:i]
				} else {
					s = string(append([]byte(s[:i]), s[i+1:]...))
				}
				i--
			}
		}
	}
	// 遍历完后如果栈不为空，说明存在'('没有匹配，故依次pop栈顶元素，逐个删除。
	for len(stack) > 0 {
		i := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		// 注意：事实上，s[i+1:]不会出现索引越界
		if i == len(s)-1 {
			s = s[:i]
		} else {
			s = string(append([]byte(s[:i]), s[i+1:]...))
		}
	}
	return s
}

// minRemoveToMakeValid_2 栈(优化 删除字符串元素的步骤)
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路：
//	a.直接遍历s，遇到'('，将他的位置入栈；遇到')'有两种可能:
//		栈里有'('：就出栈;
//		没有'(': 即栈为空直接删除当前符号，此时s长度和下标需要--。
//	b.遍历完后如果栈不为空，说明存在'('没有匹配，故依次pop栈顶元素，逐个删除。
func minRemoveToMakeValid_2(s string) string {
	var stack []int
	for i := 0; i < len(s); i++ {
		// 直接遍历s，遇到'('，将他的位置入栈
		if s[i] == '(' {
			stack = append(stack, i)
		}
		// 遇到')'有两种可能
		if s[i] == ')' {
			if len(stack) > 0 {
				// 栈里有'('：就出栈
				stack = stack[:len(stack)-1]
			} else {
				// 没有'(': 即栈为空直接删除当前符号
				// 注意：s[i+1:]不会出现索引越界
				s = string(append([]byte(s[:i]), s[i+1:]...))
				i--
			}
		}
	}
	// 遍历完后如果栈不为空，说明存在'('没有匹配，故依次pop栈顶元素，逐个删除。
	for len(stack) > 0 {
		i := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		// 注意：s[i+1:]不会出现索引越界
		s = string(append([]byte(s[:i]), s[i+1:]...))
	}
	return s
}
