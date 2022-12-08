package leetcode_0669_trim_a_binary_search_tree

// 669. 修剪二叉搜索树
// https://leetcode.cn/problems/trim-a-binary-search-tree/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// trimBST dfs
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 注意: 这是一颗二叉搜索树
func trimBST(root *TreeNode, low int, high int) *TreeNode {
	if root == nil {
		return nil
	}

	if root.Val < low {
		return trimBST(root.Right, low, high)
	}

	if root.Val > high {
		return trimBST(root.Left, low, high)
	}
	root.Left = trimBST(root.Left, low, high)
	root.Right = trimBST(root.Right, low, high)
	return root
}
