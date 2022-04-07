package leetcode_0199_binary_tree_right_side_view

// 0199.二叉树的右视图
// https://leetcode-cn.com/problems/binary-tree-right-side-view/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// rightSideView 迭代
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func rightSideView(root *TreeNode) []int {
	var (
		ans   []int
		stack []*TreeNode
	)
	if root == nil {
		return ans
	}
	stack = append(stack, root)
	for len(stack) > 0 {
		size := len(stack)
		for i := 0; i < size; i++ {
			node := stack[0]
			stack = stack[1:]
			if node.Left != nil {
				stack = append(stack, node.Left)
			}
			if node.Right != nil {
				stack = append(stack, node.Right)
			}
			if i == size-1 {
				ans = append(ans, node.Val)
			}
		}
	}
	return ans
}
