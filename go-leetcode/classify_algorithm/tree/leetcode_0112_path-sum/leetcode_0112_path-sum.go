package leetcode_0112_path_sum

// 0112.路径总和
// https://leetcode-cn.com/problems/path-sum/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// hasPathSum 递归
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 注意: 路径总和sum没有说一定是正的；可以为负数，所以必须遍历到叶子节点才能结束，不能提前结束。
func hasPathSum(root *TreeNode, targetSum int) bool {
	// 路径总和sum没有说一定是正的；可以为负数，所以必须遍历到叶子节点才能结束，不能提前结束。
	// 条件不能是  root == nil || targetSum < 0
	if root == nil {
		return false
	}
	if root.Left == nil && root.Right == nil {
		return targetSum-root.Val == 0
	}
	return hasPathSum(root.Left, targetSum-root.Val) || hasPathSum(root.Right, targetSum-root.Val)
}

// hasPathSum_2 迭代
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 注意: 路径总和sum没有说一定是正的；可以为负数，所以必须遍历到叶子节点才能结束，不能提前结束。
func hasPathSum_2(root *TreeNode, targetSum int) bool {
	// 路径总和sum没有说一定是正的；可以为负数，所以必须遍历到叶子节点才能结束，不能提前结束。
	// 条件不能是  root == nil || targetSum < 0
	if root == nil {
		return false
	}
	var (
		qNode []*TreeNode
		qVal  []int
	)
	qNode = append(qNode, root)
	qVal = append(qVal, root.Val)
	for len(qNode) > 0 {
		now := qNode[0]
		sum := qVal[0]
		qNode = qNode[1:]
		qVal = qVal[1:]
		if now.Left == nil && now.Right == nil {
			if sum == targetSum {
				return true
			}
			continue
		}
		if now.Left != nil {
			qNode = append(qNode, now.Left)
			qVal = append(qVal, now.Left.Val+sum)
		}
		if now.Right != nil {
			qNode = append(qNode, now.Right)
			qVal = append(qVal, now.Right.Val+sum)
		}
	}
	return false
}
