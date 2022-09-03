package leetcode_0776_split_bst

// 0776. 拆分二叉搜索树
// https://leetcode.cn/problems/split-bst/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// splitBST 递归
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func splitBST(root *TreeNode, target int) []*TreeNode {
	if root == nil {
		return []*TreeNode{nil, nil}
	}

	if root.Val <= target {
		bns := splitBST(root.Right, target)
		root.Right = bns[0]
		bns[0] = root
		return bns
	} else {
		bns := splitBST(root.Left, target)
		root.Left = bns[1]
		bns[1] = root
		return bns
	}
}
