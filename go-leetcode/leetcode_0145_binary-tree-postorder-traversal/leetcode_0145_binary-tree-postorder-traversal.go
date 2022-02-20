package leetcode_0145_binary_tree_postorder_traversal

// 二叉树的后序遍历
// https://leetcode-cn.com/problems/binary-tree-postorder-traversal/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// postorderTraversal_1 递归
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func postorderTraversal_1(root *TreeNode) []int {
	var result []int
	var dfs func(root *TreeNode)
	dfs = func(root *TreeNode) {
		if root == nil {
			return
		}
		dfs(root.Left)
		dfs(root.Right)
		result = append(result, root.Val)
	}
	dfs(root)
	return result
}

// postorderTraversal_2 迭代
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func postorderTraversal_2(root *TreeNode) []int {
	var stack []*TreeNode
	var result []int
	var prev *TreeNode
	for len(stack) > 0 || root != nil {
		for root != nil {
			stack = append(stack, root)
			root = root.Left
		}
		root = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if root.Right == nil || root.Right == prev {
			result = append(result, root.Val)
			prev = root
			root = nil
		} else {
			stack = append(stack, root)
			root = root.Right
		}
	}
	return result
}
