package leetcode_0124_binary_tree_maximum_path_sum

import "math"

// 124. 二叉树中的最大路径和🌟
// https://leetcode.cn/problems/binary-tree-maximum-path-sum/

// 同以下题目类似:
// 687. 最长同值路径
// https://leetcode.cn/problems/longest-univalue-path/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// maxPathSum
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路:
// 对于任意一个节点, 如果最大和路径包含该节点, 那么只可能是两种情况:
// 1.其左右子树中所构成的和路径值较大的那个加上该节点的值后向父节点回溯构成最大路径
// 2.左右子树都在最大路径中, 加上该节点的值构成了最终的最大路径
// 特殊示例：root = [-3], 输出结果-3
func maxPathSum(root *TreeNode) int {
	var (
		dfs func(root *TreeNode) int
		ans = math.MinInt32
	)

	// dfs 定义从root节点为起始点出发的最大路径和
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

// maxPathSum_2
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路:
// 对于任意一个节点, 如果最大和路径包含该节点, 那么只可能是两种情况:
// 1.其左右子树中所构成的和路径值较大的那个加上该节点的值后向父节点回溯构成最大路径
// 2.左右子树都在最大路径中, 加上该节点的值构成了最终的最大路径
// 特殊示例：root = [-3], 输出结果-3
func maxPathSum_2(root *TreeNode) int {
	var (
		dfs func(root *TreeNode) int
		ans = math.MinInt32
	)

	// dfs 定义从root节点为起始点出发的最大路径和
	dfs = func(root *TreeNode) int {
		if root == nil {
			return 0
		}
		left := dfs(root.Left)
		right := dfs(root.Right)
		// 可能情况：左右都不选；只选左边；只选右边；同时选左右
		ans = max(ans, max(max(max(root.Val, root.Val+left), root.Val+right), root.Val+left+right))
		return max(root.Val, max(left+root.Val, right+root.Val))
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
