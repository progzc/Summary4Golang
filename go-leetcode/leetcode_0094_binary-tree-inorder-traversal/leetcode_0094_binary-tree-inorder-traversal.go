package leetcode_0094_binary_tree_inorder_traversal

// 二叉树的中序遍历
// https://leetcode-cn.com/problems/binary-tree-inorder-traversal/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// inorderTraversal_1 递归
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func inorderTraversal_1(root *TreeNode) []int {
	var result []int
	var dfs func(root *TreeNode)
	dfs = func(root *TreeNode) {
		if root == nil {
			return
		}
		dfs(root.Left)
		result = append(result, root.Val)
		dfs(root.Right)
	}
	dfs(root)
	return result
}

// inorderTraversal_2 中序遍历
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func inorderTraversal_2(root *TreeNode) []int {
	var stack []*TreeNode
	var result []int
	for len(stack) > 0 || root != nil {
		for root != nil {
			stack = append(stack, root)
			root = root.Left
		}
		root = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		result = append(result, root.Val)
		root = root.Right
	}
	return result
}
