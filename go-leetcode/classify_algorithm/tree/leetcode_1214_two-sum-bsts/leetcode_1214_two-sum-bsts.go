package leetcode_1214_two_sum_bsts

// 1214. 查找两棵二叉搜索树之和
// https://leetcode.cn/problems/two-sum-bsts/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// twoSumBSTs dfs
func twoSumBSTs(root1 *TreeNode, root2 *TreeNode, target int) bool {
	if root1 == nil || root2 == nil {
		return false
	}

	if root1.Val+root2.Val == target {
		return true
	}
	if root1.Val+root2.Val > target {
		return twoSumBSTs(root1.Left, root2, target) || twoSumBSTs(root1, root2.Left, target)
	}
	return twoSumBSTs(root1.Right, root2, target) || twoSumBSTs(root1, root2.Right, target)
}
