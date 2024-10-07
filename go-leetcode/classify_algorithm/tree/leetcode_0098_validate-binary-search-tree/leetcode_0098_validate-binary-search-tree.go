package leetcode_0098_validate_binary_search_tree

import "math"

// 0098.验证二叉搜索树
// https://leetcode-cn.com/problems/validate-binary-search-tree/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// isValidBST_1 递归 深度优先遍历
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func isValidBST_1(root *TreeNode) bool {
	var dfs func(root *TreeNode, lower, upper int) bool
	dfs = func(root *TreeNode, lower, upper int) bool {
		if root == nil {
			return true
		}
		if root.Val <= lower || root.Val >= upper {
			return false
		}
		return dfs(root.Left, lower, root.Val) && dfs(root.Right, root.Val, upper)
	}
	return dfs(root, math.MinInt64, math.MaxInt64)
}

// isValidBST_2 迭代 广度优先遍历
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func isValidBST_2(root *TreeNode) bool {
	var stack []*TreeNode
	pre := math.MinInt64
	for root != nil || len(stack) > 0 {
		for root != nil {
			stack = append(stack, root)
			root = root.Left
		}
		root = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if root.Val <= pre {
			return false
		}
		pre = root.Val
		root = root.Right
	}
	return true
}

// isValidBST_3 递归 深度优先遍历
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func isValidBST_3(root *TreeNode) bool {
	var dfs func(root *TreeNode)
	var ans []int
	dfs = func(root *TreeNode) {
		if root == nil {
			return
		}
		dfs(root.Left)
		ans = append(ans, root.Val)
		dfs(root.Right)
	}
	dfs(root)
	for i, num := range ans {
		if i > 0 && ans[i-1] >= num {
			return false
		}
	}
	return true
}
