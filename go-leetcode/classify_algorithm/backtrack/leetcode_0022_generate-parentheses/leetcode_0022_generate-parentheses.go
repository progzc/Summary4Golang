package leetcode_0022_generate_parentheses

// 0022.括号生成
// https://leetcode-cn.com/problems/generate-parentheses/

// generateParenthesis 深度优先遍历（回溯）
// 思路：回溯+剪枝
//	a.可以生出左括号的条件：左括号的剩余数量大于0
//	b.可以生出右括号的条件：左括号的剩余数量小于右括号的剩余数量
func generateParenthesis(n int) []string {
	var (
		ans []string
		dfs func(left, right int, curStr string)
	)
	if n == 0 {
		return ans
	}
	// left 剩余可用的左括号
	// right 剩余可用的右括号
	dfs = func(left, right int, curStr string) {
		// 因为每一次尝试，都使用新的字符串变量，所以无需回溯
		// 亦即递归终止条件
		if left == 0 && right == 0 {
			ans = append(ans, curStr)
		}
		// 剪枝
		if left > right {
			return
		}
		// 添加左括号或右括号
		if left > 0 {
			dfs(left-1, right, curStr+"(")
		}

		if right > 0 {
			dfs(left, right-1, curStr+")")
		}
	}
	dfs(n, n, "")
	return ans
}
