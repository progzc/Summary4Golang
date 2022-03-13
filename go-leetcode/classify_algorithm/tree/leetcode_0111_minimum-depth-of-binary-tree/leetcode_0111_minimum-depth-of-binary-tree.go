package leetcode_0111_minimum_depth_of_binary_tree

import (
	"math"
)

// 0111.二叉树的最小深度
// https://leetcode-cn.com/problems/minimum-depth-of-binary-tree/

// 说明:
// 最小深度是从根节点到最近叶子节点的最短路径上的节点数量。
// 叶子节点是指没有子节点的节点。

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// minDepth_1 递归
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func minDepth_1(root *TreeNode) int {
	if root == nil {
		return 0
	}
	if root.Left == nil {
		return minDepth_1(root.Right) + 1
	}

	if root.Right == nil {
		return minDepth_1(root.Left) + 1
	}

	return min(minDepth_1(root.Left), minDepth_1(root.Right)) + 1
}

// minDepth_2 递归
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func minDepth_2(root *TreeNode) int {
	if root == nil {
		return 0
	}
	if root.Left == nil && root.Right == nil {
		return 1
	}

	depth := math.MaxInt64
	if root.Left != nil {
		depth = min(minDepth_2(root.Left), depth)
	}

	if root.Right != nil {
		depth = min(minDepth_2(root.Right), depth)
	}

	return depth + 1
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
