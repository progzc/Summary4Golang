package leetcode_0144_binary_tree_preorder_traversal

// 二叉树的前序遍历
// https://leetcode-cn.com/problems/binary-tree-preorder-traversal/
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// preorderTraversal_1 递归解法
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func preorderTraversal_1(root *TreeNode) []int {
	var result []int
	var dfs func(root *TreeNode)
	dfs = func(root *TreeNode) {
		if root == nil {
			return
		}
		result = append(result, root.Val)
		dfs(root.Left)
		dfs(root.Right)
	}
	dfs(root)
	return result
}

// preorderTraversal_2 迭代
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func preorderTraversal_2(root *TreeNode) []int {
	var stack []*TreeNode
	var result []int
	for len(stack) > 0 || root != nil {
		for root != nil {
			result = append(result, root.Val)
			stack = append(stack, root)
			root = root.Left
		}
		root = stack[len(stack)-1].Right
		stack = stack[:len(stack)-1]
	}
	return result
}
