package leetcode_0032_longest_valid_parentheses

// 0032.最长有效括号
// https://leetcode-cn.com/problems/longest-valid-parentheses/

// longestValidParentheses 栈
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路：保持栈底元素为当前已经遍历过的元素中「最后一个没有被匹配的右括号的下标」，这样的做法主要是考虑了边界条件的处理，栈里其他元素维护左括号的下标。
//	对于遇到的每个'('，我们将它的下标放入栈中
//	对于遇到的每个')'，我们先弹出栈顶元素表示匹配了当前右括号
func longestValidParentheses(s string) int {
	var (
		maxLen int
		stack  []int
	)
	n := len(s)
	stack = append(stack, -1)
	for i := 0; i < n; i++ {
		if s[i] == '(' {
			stack = append(stack, i)
		} else {
			stack = stack[:len(stack)-1]
			if len(stack) == 0 {
				stack = append(stack, i)
			} else {
				maxLen = max(maxLen, i-stack[len(stack)-1])
			}
		}
	}
	return maxLen
}

// longestValidParentheses_2 动态规划
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 注意：这个题与下面题不一样，状态转移方程不一样
//	0005.最长回文子串
//	https://leetcode-cn.com/problems/longest-palindromic-substring/
// 例如：
//	输入: s=")()())"
//	输出: 4
// 思路：
//	状态: dp[i]表示以下表i结尾的最长有效括号的长度。
//	转移方程:
//		当s[i] == ')' && s[i-1] = '('：
//			dp[i] = dp[i-2]+2
//		当s[i] == ')' && s[i-1] = ')' && s[i-dp[i-1]-1] = '('：
//			dp[i] = dp[i-1]+2+dp[i-dp[i-1]-2]
func longestValidParentheses_2(s string) int {
	n := len(s)
	maxLen := 0
	dp := make([]int, n)
	for i := 1; i < n; i++ {
		if s[i] == '(' {
			continue
		}
		if s[i-1] == '(' {
			if i >= 2 {
				dp[i] = dp[i-2] + 2
			} else {
				dp[i] = 2
			}
		} else if i-dp[i-1]-1 >= 0 && s[i-dp[i-1]-1] == '(' {
			if i-dp[i-1]-2 >= 0 {
				dp[i] = dp[i-1] + dp[i-dp[i-1]-2] + 2
			} else {
				dp[i] = dp[i-1] + 2
			}
		}
		maxLen = max(maxLen, dp[i])
	}
	return maxLen
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
