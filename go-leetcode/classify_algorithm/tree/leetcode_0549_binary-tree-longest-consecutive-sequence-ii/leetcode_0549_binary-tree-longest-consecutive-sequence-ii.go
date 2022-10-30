package leetcode_0549_binary_tree_longest_consecutive_sequence_ii

// 0549. 二叉树中最长的连续序列
// https://leetcode.cn/problems/binary-tree-longest-consecutive-sequence-ii/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// longestConsecutive dfs
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路: 在每一个点，我们使用两个变量 inr 和 dcr，其中 inr 表示当前点为止最长增长序列的长度（包括该点自己），
//	dcr 表示当前点为止最长下降序列的长度（包括该点自己）。
// 注意: 要是整个链路是上升或下降的，且相邻节点值差值为1
// 例1:
//	输入: [4,-7,-3,null,null,-9,-3,9,-7,-4,null,6,null,-6,-6,null,null,0,6,5,null,9,null,null,-1,-4,null,null,null,-2]
//	输出: 2
func longestConsecutive(root *TreeNode) int {
	var (
		dfs func(root *TreeNode) (int, int)
		ans int
	)
	// 在每一个点，我们使用两个变量 inr 和 dcr，其中 inr 表示当前点为止最长增长序列的长度（包括该点自己），
	// dcr 表示当前点为止最长下降序列的长度（包括该点自己）
	dfs = func(root *TreeNode) (int, int) {
		if root == nil {
			return 0, 0
		}
		var (
			incr = 1
			dcr  = 1
		)
		if root.Left != nil {
			incr1, dcr1 := dfs(root.Left)
			if root.Val == root.Left.Val+1 {
				dcr = dcr1 + 1
			} else if root.Val == root.Left.Val-1 {
				incr = incr1 + 1
			}
		}
		if root.Right != nil {
			incr2, dcr2 := dfs(root.Right)
			if root.Val == root.Right.Val+1 {
				dcr = max(dcr, dcr2+1)
			} else if root.Val == root.Right.Val-1 {
				incr = max(incr, incr2+1)
			}
		}
		ans = max(ans, incr+dcr-1)
		return incr, dcr
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
