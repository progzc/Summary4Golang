package leetcode_0687_longest_univalue_path

// 687. 最长同值路径
// https://leetcode.cn/problems/longest-univalue-path/

// 同以下题目类似:
// 124. 二叉树中的最大路径和
// https://leetcode.cn/problems/binary-tree-maximum-path-sum/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// longestUnivaluePath
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路: 最长同值路径长度必定为某一节点的左最长同值有向路径长度与右最长同值有向路径长度之和。
// 注意: 两个节点之间的路径长度 由它们之间的边数表示。
func longestUnivaluePath(root *TreeNode) int {
	var (
		ans int
		dfs func(root *TreeNode) int
	)

	// dfs 定义从root节点为起始点出发的最长同值路径长度
	dfs = func(root *TreeNode) int {
		if root == nil {
			return 0
		}
		// 计算左节点的最长同值路径长度
		left := dfs(root.Left)
		// 计算右节点的最长同值路径长度
		right := dfs(root.Right)
		if root.Left != nil && root.Val == root.Left.Val {
			left++
		} else {
			left = 0
		}
		if root.Right != nil && root.Val == root.Right.Val {
			right++
		} else {
			right = 0
		}
		ans = max(ans, left+right)
		return max(left, right)
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
