package leetcode_0515_find_largest_value_in_each_tree_row

import "math"

// 515. 在每个树行中找最大值
// https://leetcode.cn/problems/find-largest-value-in-each-tree-row/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// largestValues 迭代
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func largestValues(root *TreeNode) []int {
	var (
		ans   []int
		stack []*TreeNode
	)

	if root == nil {
		return ans
	}
	stack = append(stack, root)
	for len(stack) > 0 {
		n := len(stack)
		curMax := math.MinInt32
		for i := 0; i < n; i++ {
			node := stack[0]
			stack = stack[1:]
			if node.Val > curMax {
				curMax = node.Val
			}
			if node.Left != nil {
				stack = append(stack, node.Left)
			}
			if node.Right != nil {
				stack = append(stack, node.Right)
			}
		}
		ans = append(ans, curMax)
	}
	return ans
}
