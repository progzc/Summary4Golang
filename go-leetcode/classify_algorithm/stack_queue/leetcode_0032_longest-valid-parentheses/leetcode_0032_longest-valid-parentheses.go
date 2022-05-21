package leetcode_0032_longest_valid_parentheses

import "strings"

// 0032.最长有效括号
// https://leetcode-cn.com/problems/longest-valid-parentheses/

// longestValidParentheses 栈
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路：保持栈底元素为当前已经遍历过的元素中「最后一个没有被匹配的右括号的下标」，这样的做法主要是考虑了边界条件的处理，栈里其他元素维护左括号的下标。
//	对于遇到的每个'('，我们将它的下标放入栈中
//	对于遇到的每个')'，我们先弹出栈顶元素表示匹配了当前右括号
func longestValidParentheses(s string) int {

}
