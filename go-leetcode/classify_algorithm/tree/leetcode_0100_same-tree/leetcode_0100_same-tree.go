package leetcode_0100_same_tree

// 0100.相同的树
// https://leetcode-cn.com/problems/same-tree/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// isSameTree_1 递归
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func isSameTree_1(p *TreeNode, q *TreeNode) bool {
	if p == nil && q == nil {
		return true
	}
	if p == nil || q == nil {
		return false
	}
	if p.Val != q.Val {
		return false
	}
	return isSameTree_1(p.Left, q.Left) && isSameTree_1(p.Right, q.Right)
}
