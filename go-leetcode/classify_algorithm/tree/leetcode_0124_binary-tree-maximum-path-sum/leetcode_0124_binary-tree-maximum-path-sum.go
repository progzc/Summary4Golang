package leetcode_0124_binary_tree_maximum_path_sum

import "math"

// 124. 二叉树中的最大路径和
// https://leetcode.cn/problems/binary-tree-maximum-path-sum/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// maxPathSum
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路:
//	对于任意一个节点, 如果最大和路径包含该节点, 那么只可能是两种情况:
//	  1.其左右子树中所构成的和路径值较大的那个加上该节点的值后向父节点回溯构成最大路径
//	  2.左右子树都在最大路径中, 加上该节点的值构成了最终的最大路径
func maxPathSum(root *TreeNode) int {
	var (
		dfs func(root *TreeNode) int
		ans = math.MinInt32
	)
	dfs = func(root *TreeNode) int {
		if root == nil {
			return 0
		}
		// 计算左边分支最大值，左边分支如果为负数还不如不选择
		left := max(dfs(root.Left), 0)
		// 计算右边分支最大值，右边分支如果为负数还不如不选择
		right := max(dfs(root.Right), 0)
		// left->root->right 作为路径与已经计算过历史最大值做比较
		ans = max(ans, root.Val+left+right)
		// 返回经过root的单边最大分支给当前root的父节点计算使用
		return root.Val + max(left, right)
	}
	dfs(root)
	return ans
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
