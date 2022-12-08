package leetcode_0655_print_binary_tree

import (
	"strconv"
)

// 655. 输出二叉树
// https://leetcode.cn/problems/print-binary-tree/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// printTree dfs+二叉树节点编号
// 时间复杂度: O(height*2^height)
// 空间复杂度: O(height)
func printTree(root *TreeNode) [][]string {
	var (
		ans      [][]string
		calDepth func(root *TreeNode) int
		dfs      func(root *TreeNode, r, c int)
	)

	if root == nil {
		return ans
	}

	// 1.计算树的高度
	calDepth = func(root *TreeNode) int {
		if root == nil {
			return 0
		}
		return max(calDepth(root.Left)+1, calDepth(root.Right)+1)
	}
	height := calDepth(root) - 1
	// 2.初始化结果
	m, n := height+1, 1<<(height+1)-1
	ans = make([][]string, m)
	for i := 0; i < m; i++ {
		ans[i] = make([]string, n)
	}
	// 3.dfs
	dfs = func(root *TreeNode, r, c int) {
		ans[r][c] = strconv.Itoa(root.Val)
		if root.Left != nil {
			dfs(root.Left, r+1, c-1<<(height-r-1))
		}
		if root.Right != nil {
			dfs(root.Right, r+1, c+1<<(height-r-1))
		}
	}
	dfs(root, 0, (n-1)/2)
	return ans
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
