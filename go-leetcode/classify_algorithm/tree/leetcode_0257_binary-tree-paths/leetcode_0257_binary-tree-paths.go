package leetcode_0257_binary_tree_paths

import (
	"fmt"
	"strings"
)

// 257. 二叉树的所有路径
// https://leetcode.cn/problems/binary-tree-paths/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// binaryTreePaths
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func binaryTreePaths(root *TreeNode) []string {
	var (
		ans  []string
		path []string
		dfs  func(root *TreeNode)
	)
	dfs = func(root *TreeNode) {
		if root == nil {
			return
		}
		path = append(path, fmt.Sprintf("%d", root.Val))
		defer func() {
			path = path[:len(path)-1]
		}()
		if root.Left == nil && root.Right == nil {
			ans = append(ans, strings.Join(append([]string(nil), path...), "->"))
			return
		}
		dfs(root.Left)
		dfs(root.Right)
	}
	dfs(root)
	return ans
}
