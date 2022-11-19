package leetcode_0144_binary_tree_preorder_traversal

// 144. 二叉树的前序遍历
// https://leetcode.cn/problems/binary-tree-preorder-traversal/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// preorderTraversal 递归法
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func preorderTraversal(root *TreeNode) []int {
	var (
		ans []int
		dfs func(root *TreeNode)
	)
	dfs = func(root *TreeNode) {
		if root == nil {
			return
		}
		ans = append(ans, root.Val)
		dfs(root.Left)
		dfs(root.Right)
	}
	dfs(root)
	return ans
}

// preorderTraversal_2 迭代法
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func preorderTraversal_2(root *TreeNode) []int {
	var (
		ans   []int
		stack []*TreeNode
	)
	if root == nil {
		return ans
	}
	stack = append(stack, root)
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		ans = append(ans, node.Val)
		if node.Right != nil {
			stack = append(stack, node.Right)
		}
		if node.Left != nil {
			stack = append(stack, node.Left)
		}
	}
	return ans
}
