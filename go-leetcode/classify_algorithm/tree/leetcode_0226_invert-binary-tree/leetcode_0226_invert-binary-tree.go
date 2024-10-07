package leetcode_0226_invert_binary_tree

// 0226. 翻转二叉树
// https://leetcode.cn/problems/invert-binary-tree

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// invertTree 递归
// 时间复杂度: O(n)
// 空间复杂度: O(log(n))
func invertTree(root *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	right := invertTree(root.Right)
	left := invertTree(root.Left)
	root.Left = right
	root.Right = left
	return root
}
