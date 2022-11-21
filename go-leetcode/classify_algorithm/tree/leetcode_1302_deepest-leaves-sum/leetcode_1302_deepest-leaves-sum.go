package leetcode_1302_deepest_leaves_sum

// 1302. 层数最深叶子节点的和
// https://leetcode.cn/problems/deepest-leaves-sum/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// deepestLeavesSum 迭代
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func deepestLeavesSum(root *TreeNode) int {
	var (
		sum   int
		stack []*TreeNode
	)
	if root == nil {
		return sum
	}
	stack = append(stack, root)
	for len(stack) > 0 {
		var (
			ss int
			n  = len(stack)
		)
		for i := 0; i < n; i++ {
			node := stack[0]
			stack = stack[1:]
			ss += node.Val
			if node.Left != nil {
				stack = append(stack, node.Left)
			}
			if node.Right != nil {
				stack = append(stack, node.Right)
			}
		}
		sum = ss
	}
	return sum
}
