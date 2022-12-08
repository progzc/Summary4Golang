package leetcode_0654_maximum_binary_tree

import "math"

// 654. 最大二叉树
// https://leetcode.cn/problems/maximum-binary-tree/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// constructMaximumBinaryTree dfs
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func constructMaximumBinaryTree(nums []int) *TreeNode {
	if len(nums) == 0 {
		return nil
	}

	cur := math.MinInt32
	idx := -1
	for i, num := range nums {
		if num > cur {
			idx = i
			cur = num
		}
	}

	root := &TreeNode{Val: cur}
	root.Left = constructMaximumBinaryTree(nums[0:idx])
	root.Right = constructMaximumBinaryTree(nums[idx+1:])
	return root
}
