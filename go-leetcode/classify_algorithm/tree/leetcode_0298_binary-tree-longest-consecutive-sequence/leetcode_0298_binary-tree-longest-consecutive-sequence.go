package leetcode_0298_binary_tree_longest_consecutive_sequence

// 0298. 二叉树最长连续序列
// https://leetcode.cn/problems/binary-tree-longest-consecutive-sequence/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// longestConsecutive 自顶向下深度优先搜索
// 时间复杂度: O(n)
// 空间复杂度: (n)
func longestConsecutive(root *TreeNode) int {
	var (
		maxLen int
		dfs    func(cur, parent *TreeNode, length int)
	)
	dfs = func(cur, parent *TreeNode, length int) {
		if cur == nil {
			return
		}

		if parent != nil && cur.Val == parent.Val+1 {
			length += 1
		} else {
			length = 1
		}

		maxLen = max(maxLen, length)
		dfs(cur.Left, cur, length)
		dfs(cur.Right, cur, length)
	}
	dfs(root, nil, 0)
	return maxLen
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
