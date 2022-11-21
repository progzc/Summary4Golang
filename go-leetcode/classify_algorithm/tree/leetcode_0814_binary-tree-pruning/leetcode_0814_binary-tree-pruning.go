package leetcode_0814_binary_tree_pruning

// 814. 二叉树剪枝
// https://leetcode.cn/problems/binary-tree-pruning/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// pruneTree 递归
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func pruneTree(root *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	left := pruneTree(root.Left)
	right := pruneTree(root.Right)
	if root.Val == 0 && left == nil && right == nil {
		return nil
	} else {
		root.Left = left
		root.Right = right
	}
	return root
}
