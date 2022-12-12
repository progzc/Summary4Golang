package leetcode_0678_valid_parenthesis_string

// 678. 有效的括号字符串
// https://leetcode.cn/problems/valid-parenthesis-string/

// 类似题目:
//	[20.有效的括号](https://leetcode.cn/problems/valid-parentheses/)
//	[32.最长有效括号](https://leetcode.cn/problems/longest-valid-parentheses/)

// checkValidString
// 时间复杂度: O(n^3)
// 空间复杂度: O(n^2)
// 思路: 要判断 s 是否为有效的括号字符串，需要判断 s 的首尾字符以及 s 的中间字符是否符合有效的括号字符串的要求。可以使用动态规划求解。
//	状态: dp[i][j]表示字符串 s 从下标 i 到 j 的子串是否为有效的括号字符串。
//	转移方程: 当j-i>=2时,dp[i][j]满足下面任一条件即为true:
//		a.dp[i][j] = dp[i+1][j-1] && ( s[i] 和 s[j] 分别为左括号和右括号，或者为 '*' )
//		b.存在i<=k<j,使得dp[i][k]和dp[k+1][j]都为true, 则dp[i][j]=true。简言之,将两个子串拼接之后的子串也为有效的括号字符串。
//			典序示例: "()()"、"*)()"
//	边界条件:
//		a.当j-i==0时,只有当该字符串是"*"时，才是有效的括号字符串。
//		b.当j-i==1时,只有当该字符串是"()","(*","*)","**"中的一种情况时，才是有效的括号。
func checkValidString(s string) bool {
	n := len(s)
	if n == 0 {
		return true
	}
	dp := make([][]bool, n)
	for i := 0; i < n; i++ {
		dp[i] = make([]bool, n)
	}
	// 边界条件1: 当j-i==0时,只有当该字符串是"*"时，才是有效的括号字符串。
	for i := 0; i < n; i++ {
		if s[i] == '*' {
			dp[i][i] = true
		}
	}
	// 边界条件2: 当j-i==1时,只有当该字符串是"()","(*","*)","**"中的一种情况时，才是有效的括号。
	for i := 1; i < n; i++ {
		c1, c2 := s[i-1], s[i]
		dp[i-1][i] = (c1 == '(' || c1 == '*') && (c2 == ')' || c2 == '*')
	}

	// 注意状态转移的顺序
	for i := n - 3; i >= 0; i-- {
		c1 := s[i]
		for j := i + 2; j < n; j++ {
			c2 := s[j]
			if (c1 == '(' || c1 == '*') && (c2 == ')' || c2 == '*') {
				dp[i][j] = dp[i+1][j-1]
			}
			for k := i; k < j && !dp[i][j]; k++ {
				dp[i][j] = dp[i][k] && dp[k+1][j]
			}
		}
	}
	return dp[0][n-1]
}

// checkValidString_2
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路:
//	如果字符串中没有星号，则只需要一个栈存储左括号，在从左到右遍历字符串的过程中检查括号是否匹配。
//	在有星号的情况下，需要两个栈分别存储左括号和星号。从左到右遍历字符串，进行如下操作：
//	  a.如果遇到左括号，则将当前下标存入左括号栈。
//	  b.如果遇到星号，则将当前下标存入星号栈。
//	  c.如果遇到右括号，则需要有一个左括号或星号和右括号匹配，由于星号也可以看成右括号或者空字符串，
//		因此当前的右括号应优先和左括号匹配，没有左括号时和星号匹配。
//		i)如果左括号栈不为空，则从左括号栈弹出栈顶元素
//		ii)如果左括号栈为空且星号栈不为空，则从星号栈弹出栈顶元素
//		iii)如果左括号栈和星号栈都为空，则没有字符可以和当前的右括号匹配，返回 false。
//	  d.遍历结束之后，左括号栈和星号栈可能还有元素。为了将每个左括号匹配，需要将星号看成右括号，且每个左括号必须出现在其匹配的星号之前。
//	    当两个栈都不为空时，每次从左括号栈和星号栈分别弹出栈顶元素，对应左括号下标和星号下标，判断是否可以匹配，
//	    匹配的条件是左括号下标小于星号下标，如果左括号下标大于星号下标则返回 false。
//	  e.最终判断左括号栈是否为空。如果左括号栈为空，则左括号全部匹配完毕，剩下的星号都可以看成空字符串，
//	    此时 s 是有效的括号字符串，返回 true。如果左括号栈不为空，则还有左括号无法匹配，此时 s 不是有效的括号字符串，返回 false。
func checkValidString_2(s string) bool {
	n := len(s)
	if n == 0 {
		return true
	}
	var leftStack, starStack []int
	for i, ch := range s {
		if ch == '(' {
			leftStack = append(leftStack, i)
		} else if ch == '*' {
			starStack = append(starStack, i)
		} else {
			if len(leftStack) > 0 {
				leftStack = leftStack[:len(leftStack)-1]
			} else if len(starStack) > 0 {
				starStack = starStack[:len(starStack)-1]
			} else {
				return false
			}
		}
	}
	for len(leftStack) > 0 {
		if len(starStack) > 0 && leftStack[len(leftStack)-1] < starStack[len(starStack)-1] {
			leftStack = leftStack[:len(leftStack)-1]
			starStack = starStack[:len(starStack)-1]
		} else {
			return false
		}
	}
	return true
}
